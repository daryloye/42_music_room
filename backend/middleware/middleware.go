package middleware

import (
	"context"
	"net/http"
	"server/config"
	"server/helper"

	"github.com/julienschmidt/httprouter"
)

var userIdKey = "userId"

type Middleware interface {
	RequireAuth(next httprouter.Handle) httprouter.Handle
}

type MiddlewareImpl struct {
	Cfg *config.Config
}

func NewMiddleware(cfg *config.Config) Middleware {
	return &MiddlewareImpl{Cfg: cfg}
}

func (m *MiddlewareImpl) RequireAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		accessToken, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, helper.ErrUserInvalidOrExpiredToken.Error(), http.StatusUnauthorized)
			return
		}

		userId, err := helper.DecodeJWT(accessToken.Value, m.Cfg.EnvJwtAccessSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userId)

		next(w, r.WithContext(ctx), params)
	}
}

func UserIdFromContext(ctx context.Context) string {
	userId := ctx.Value(userIdKey).(string)
	return userId
}
