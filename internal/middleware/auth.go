package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mikkkkkkka/what-i-know-api/internal/api"
	"github.com/mikkkkkkka/what-i-know-api/internal/auth"
	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type ctxKey string

var CtxUserID ctxKey = "user_id"

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

func (m *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.WriteError(w, domain.ErrIncorrectCredentials)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.jwtManager.ParseJWTToken(tokenString)
		if err != nil {
			api.WriteError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
