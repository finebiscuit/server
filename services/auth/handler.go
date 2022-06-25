package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	authv1 "github.com/finebiscuit/proto/biscuit/auth/v1"
	"github.com/finebiscuit/proto/biscuit/auth/v1/authv1connect"
)

func NewHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	h := &handler{}
	return authv1connect.NewAuthHandler(h, opts...)
}

type handler struct {
}

func (h *handler) SignUp(
	ctx context.Context,
	req *connect.Request[authv1.SignUpRequest],
) (*connect.Response[authv1.SignUpResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (h *handler) CreateSession(
	ctx context.Context,
	req *connect.Request[authv1.CreateSessionRequest],
) (*connect.Response[authv1.CreateSessionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}

func (h *handler) GetAccessToken(
	ctx context.Context,
	req *connect.Request[authv1.GetAccessTokenRequest],
) (*connect.Response[authv1.GetAccessTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("not implemented"))
}
