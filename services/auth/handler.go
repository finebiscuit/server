package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	authv1 "github.com/finebiscuit/proto/biscuit/auth/v1"
	"github.com/finebiscuit/proto/biscuit/auth/v1/authv1connect"
	"github.com/finebiscuit/server/services/auth/session"
	"github.com/finebiscuit/server/services/auth/workspace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewHandler(service Service, opts ...connect.HandlerOption) (string, http.Handler) {
	h := &handler{Auth: service}
	return authv1connect.NewAuthHandler(h, opts...)
}

type handler struct {
	Auth Service
}

func (h *handler) SignUp(
	ctx context.Context,
	req *connect.Request[authv1.SignUpRequest],
) (*connect.Response[authv1.SignUpResponse], error) {
	ident, err := h.Auth.SignUp(ctx, req.Msg.GetEmail(), req.Msg.GetPassword())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&authv1.SignUpResponse{
		Identity: &authv1.Identity{
			UserId:    ident.UserID,
			UserEmail: ident.UserEmail,
		},
	})
	return res, nil
}

func (h *handler) CreateSession(
	ctx context.Context,
	req *connect.Request[authv1.CreateSessionRequest],
) (*connect.Response[authv1.CreateSessionResponse], error) {
	sess, ident, ws, err := h.Auth.CreateSession(ctx, req.Msg.GetLogin(), req.Msg.GetPassword())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&authv1.CreateSessionResponse{
		Session: &authv1.Session{
			Id:        sess.ID.String(),
			Code:      sess.PlainCode,
			CreatedAt: timestamppb.New(sess.CreatedAt),
			ExpiresAt: timestamppb.New(sess.ExpiresAt),
		},
		Identity: &authv1.Identity{
			UserId:    ident.UserID,
			UserEmail: ident.UserEmail,
		},
		Workspaces: make([]*authv1.Workspace, 0, len(ws)),
	})

	for _, w := range ws {
		res.Msg.Workspaces = append(res.Msg.Workspaces, &authv1.Workspace{
			Id:          w.ID.String(),
			DisplayName: w.DisplayName,
		})
	}

	return res, nil
}

func (h *handler) GetAccessToken(
	ctx context.Context,
	req *connect.Request[authv1.GetAccessTokenRequest],
) (*connect.Response[authv1.GetAccessTokenResponse], error) {
	sessID, err := session.ParseID(req.Msg.GetSessionId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	wsID, err := workspace.ParseID(req.Msg.GetWorkspaceId())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	tok, signedTok, err := h.Auth.GetAccessToken(ctx, sessID, req.Msg.GetSessionCode(), wsID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	res := connect.NewResponse(&authv1.GetAccessTokenResponse{
		AccessToken: signedTok,
		ExpiresIn:   uint32(tok.GetTTL().Seconds()),
		TokenType:   "Bearer",
	})
	return res, nil
}

func (h *handler) VerifyAccessToken(
	ctx context.Context,
	req *connect.Request[authv1.VerifyAccessTokenRequest],
) (*connect.Response[authv1.VerifyAccessTokenResponse], error) {
	tok, err := h.Auth.VerifyAccessToken(ctx, req.Msg.GetAccessToken())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	res := connect.NewResponse(&authv1.VerifyAccessTokenResponse{
		UserId:      tok.UserID.String(),
		WorkspaceId: tok.WorkspaceID.String(),
		SessionId:   tok.SessionID.String(),
		ExpiresIn:   uint32(tok.ExpiresAt.Sub(time.Now()).Seconds()),
		IssuedAt:    timestamppb.New(tok.IssuedAt),
	})
	return res, nil
}
