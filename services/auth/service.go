package auth

import (
	"context"
	"time"

	"github.com/finebiscuit/server/services/auth/session"
	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
)

type Service interface {
	SignUp(ctx context.Context, email, password string) (*user.Identity, error)
	CreateSession(ctx context.Context, login, password string) (*session.Session, *user.Identity, []*workspace.Workspace, error)
	GetAccessToken(ctx context.Context, sessID session.ID, sessCode string, wsID workspace.ID) (*session.AccessToken, string, error)
	VerifyAccessToken(ctx context.Context, token string) (*session.AccessToken, error)
}

func NewService(tx TxFn, secretKey []byte) Service {
	return &serviceImpl{
		tx:        tx,
		secretKey: secretKey,
	}
}

type serviceImpl struct {
	tx        TxFn
	secretKey []byte
}

func (s *serviceImpl) SignUp(ctx context.Context, email, password string) (*user.Identity, error) {
	u, err := user.NewFromEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}

	w, err := workspace.New(u.ID, "")
	if err != nil {
		return nil, err
	}

	err = s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Users().Create(ctx, u); err != nil {
			return err
		}

		if err := uow.Workspaces().Create(ctx, w); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.AsIdentity(), nil
}

func (s *serviceImpl) CreateSession(
	ctx context.Context, login, password string,
) (*session.Session, *user.Identity, []*workspace.Workspace, error) {
	var (
		sess *session.Session
		u    *user.User
		ws   []*workspace.Workspace
	)
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) (err error) {
		u, err = uow.Users().GetByLogin(ctx, login)
		if err != nil {
			return err
		}

		if err := u.ComparePassword(password); err != nil {
			return err
		}

		sess, err = session.New(u.ID, 30*24*time.Hour)
		if err != nil {
			return err
		}

		if err := uow.Sessions().Create(ctx, sess); err != nil {
			return err
		}

		u, err = uow.Users().Get(ctx, sess.UserID)
		if err != nil {
			return err
		}

		ws, err = uow.Workspaces().List(ctx, u.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return sess, u.AsIdentity(), ws, nil
}

func (s *serviceImpl) GetAccessToken(
	ctx context.Context, sessID session.ID, sessCode string, wsID workspace.ID,
) (*session.AccessToken, string, error) {
	var (
		tok       *session.AccessToken
		signedTok string
	)
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		sess, err := uow.Sessions().Get(ctx, sessID)
		if err != nil {
			return err
		}

		if err := sess.CompareCode(sessCode); err != nil {
			return err
		}

		ws, err := uow.Workspaces().Get(ctx, wsID)
		if err != nil {
			return err
		}

		if err := ws.CompareAccessFor(sess.UserID); err != nil {
			return err
		}

		tok, err = sess.GenerateAccessToken(wsID, time.Now().Add(5*time.Minute))
		if err != nil {
			return err
		}

		signedTok, err = tok.SignedString([]byte(s.secretKey))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return tok, signedTok, nil
}

func (s *serviceImpl) VerifyAccessToken(ctx context.Context, token string) (*session.AccessToken, error) {
	accessToken, err := session.ParseAccessToken(token, []byte(s.secretKey))
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}
