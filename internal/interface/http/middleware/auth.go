package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/smetanamolokovich/veylo/pkg/jwt"
)

type contextKey string

const (
	orgIDKey  contextKey = "organization_id"
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func Auth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(tokenStr, prefix) {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}
			tokenStr = strings.TrimPrefix(tokenStr, prefix)

			claims, err := jwtManager.Validate(tokenStr)
			if err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, orgIDKey, claims.OrganizationID)
			ctx = context.WithValue(ctx, userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func OrganizationIDFromCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(orgIDKey).(string)
	return v, ok
}

func UserIDFromCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}

func RoleFromCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(roleKey).(string)
	return v, ok
}
