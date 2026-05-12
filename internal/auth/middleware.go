package auth

import (
	"context"
	"net/http"
	"strings"

	"clinic-api/internal/shared"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

type Middleware struct {
	jwtManager *JWTManager
}

func NewMiddleware(jwtManager *JWTManager) *Middleware {
	return &Middleware{jwtManager: jwtManager}
}

func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader == "" {
			shared.WriteError(w, http.StatusUnauthorized, "token requerido")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			shared.WriteError(w, http.StatusUnauthorized, "formato de token invalido")
			return
		}

		claims, err := m.jwtManager.Validate(parts[1])
		if err != nil {
			shared.WriteError(w, http.StatusUnauthorized, "token invalido o expirado")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(claimsKey).(*Claims)
	return claims, ok
}
