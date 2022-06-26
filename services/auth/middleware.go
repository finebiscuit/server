package auth

import (
	"net/http"
	"strings"

	"github.com/finebiscuit/server/services/auth/session"
	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
)

func NewMiddleware(service Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			h := strings.Split(r.Header.Get("authorization"), " ")
			if len(h) < 2 {
				http.Error(w, "Invalid Authorization", http.StatusForbidden)
				return
			}

			token, err := service.VerifyAccessToken(r.Context(), h[len(h)-1])
			if err != nil {
				http.Error(w, "Invalid Authorization", http.StatusForbidden)
			}

			ctx := r.Context()
			ctx = session.SetAccessTokenContext(ctx, token)
			ctx = user.SetContext(ctx, token.UserID)
			ctx = workspace.SetContext(ctx, token.WorkspaceID)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
