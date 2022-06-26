package user

import "context"

type ctxType string

const ctxKey = ctxType("user_id")

func FromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(ctxKey).(ID); ok {
		return id
	}
	return ID{}
}

func SetContext(ctx context.Context, id ID) context.Context {
	// Disallow overriding user ID if present
	if _, ok := ctx.Value(ctxKey).(ID); ok {
		return ctx
	}
	return context.WithValue(ctx, ctxKey, id)
}
