package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// TokenMaker interface (assuming it's defined elsewhere)
type TokenMaker interface {
	VerifyToken(token string) (interface{}, error)
}

// AuthMiddleware creates a chi middleware for authorization
func AuthMiddleware(tokenMaker TokenMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(authorizationHeaderKey)

			if len(authorizationHeader) == 0 {
				err := errors.New("authorization header is not provided")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, errorResponse(err))
				return
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, errorResponse(err))
				return
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationTypeBearer {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, errorResponse(err))
				return
			}

			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, errorResponse(err))
				return
			}

			ctx := context.WithValue(r.Context(), authorizationPayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// errorResponse function (assuming it's defined elsewhere)
func errorResponse(err error) map[string]interface{} {
	return map[string]interface{}{"error": err.Error()}
}
