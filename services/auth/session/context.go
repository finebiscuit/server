package session

import "context"

type ctxType string

const ctxKey = ctxType("session_id")

func AccessTokenFromContext(ctx context.Context) *AccessToken {
	if token, ok := ctx.Value(ctxKey).(*AccessToken); ok {
		return token
	}
	return &AccessToken{}
}

func SetAccessTokenContext(ctx context.Context, token *AccessToken) context.Context {
	// Disallow overriding access token if present
	if _, ok := ctx.Value(ctxKey).(*AccessToken); ok {
		return ctx
	}
	return context.WithValue(ctx, ctxKey, token)
}
