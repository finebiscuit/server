package auth

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	authv1 "github.com/finebiscuit/proto/biscuit/auth/v1"
	"github.com/finebiscuit/proto/biscuit/auth/v1/authv1connect"

	"github.com/finebiscuit/server/services/auth/session"
	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
)

func NewClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) Service {
	return &client{
		Auth: authv1connect.NewAuthClient(httpClient, baseURL, opts...),
	}
}

type client struct {
	Auth authv1connect.AuthClient
}

func (c *client) SignUp(ctx context.Context, email, password string) (*user.Identity, error) {
	req := connect.NewRequest(&authv1.SignUpRequest{
		Email:    email,
		Password: password,
	})

	res, err := c.Auth.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}

	ident := &user.Identity{
		UserID:    res.Msg.GetIdentity().GetUserId(),
		UserEmail: res.Msg.GetIdentity().GetUserEmail(),
	}
	return ident, nil
}

func (c *client) CreateSession(
	ctx context.Context, login, password string,
) (*session.Session, *user.Identity, []*workspace.Workspace, error) {
	req := connect.NewRequest(&authv1.CreateSessionRequest{
		Login:    login,
		Password: password,
	})

	res, err := c.Auth.CreateSession(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	sess, err := session.NewFromProto(res.Msg.GetSession())
	if err != nil {
		return nil, nil, nil, err
	}

	ident := &user.Identity{
		UserID:    res.Msg.GetIdentity().GetUserId(),
		UserEmail: res.Msg.GetIdentity().GetUserEmail(),
	}

	wss := make([]*workspace.Workspace, 0, len(res.Msg.GetWorkspaces()))
	for _, proto := range res.Msg.GetWorkspaces() {
		ws, err := workspace.NewFromProto(proto)
		if err != nil {
			return nil, nil, nil, err
		}
		wss = append(wss, ws)
	}

	return sess, ident, wss, nil
}

func (c *client) GetAccessToken(
	ctx context.Context, sessID session.ID, sessCode string, wsID workspace.ID,
) (*session.AccessToken, string, error) {
	req := connect.NewRequest(&authv1.GetAccessTokenRequest{
		SessionId:   sessID.String(),
		SessionCode: sessCode,
		WorkspaceId: wsID.String(),
	})
	res, err := c.Auth.GetAccessToken(ctx, req)
	if err != nil {
		return nil, "", err
	}

	tok, err := c.VerifyAccessToken(ctx, res.Msg.GetAccessToken())
	if err != nil {
		return nil, "", err
	}
	return tok, res.Msg.GetAccessToken(), nil
}

func (c *client) VerifyAccessToken(ctx context.Context, token string) (*session.AccessToken, error) {
	req := connect.NewRequest(&authv1.VerifyAccessTokenRequest{
		AccessToken: token,
	})
	res, err := c.Auth.VerifyAccessToken(ctx, req)
	if err != nil {
		return nil, err
	}

	sessID, err := session.ParseID(res.Msg.GetSessionId())
	if err != nil {
		return nil, err
	}
	userID, err := user.ParseID(res.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	wsID, err := workspace.ParseID(res.Msg.GetWorkspaceId())
	if err != nil {
		return nil, err
	}

	tok := &session.AccessToken{
		SessionID:   sessID,
		UserID:      userID,
		WorkspaceID: wsID,
		IssuedAt:    res.Msg.IssuedAt.AsTime(),
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(res.Msg.GetExpiresIn())),
	}
	return tok, nil
}
