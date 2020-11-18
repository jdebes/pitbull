package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jdebes/pitbull/api"
)

const (
	authHeader = "Authorization"
	authType   = "Bearer"
)

var (
	authError    = errors.New("invalid Auth header")
	expiredError = errors.New("token expired")
)

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawAuth := r.Header.Get(authHeader)
		authSplit := strings.Split(rawAuth, authType+" ")
		if len(authSplit) != 2 {
			api.Info(w, authError, http.StatusBadRequest)
			return
		}
		token := authSplit[1]

		claims, expired, err := VerifyJwt(token)
		if err != nil {
			api.Info(w, err, http.StatusBadRequest)
			return
		}

		if expired {
			api.Info(w, expiredError, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(WithUserClaims(r.Context(), claims)))
	})
}
