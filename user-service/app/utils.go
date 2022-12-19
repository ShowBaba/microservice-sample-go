package app

import (
	"fmt"
	"net/http"

	"github.com/microservice-sample-go/shared"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidateGatewayToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("gateway_signature")
			if token == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(shared.WriteError(http.StatusBadRequest, "gateway-token not in header", nil))
				if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, "gateway-token not in header"); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
					return
				}
				return
			}
			claim, err := shared.ValidateGatewayToken(token, GetConfig().UserServiceSecretKey)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(shared.WriteError(http.StatusBadRequest, "error validating gateway-token", fmt.Sprintf("%v", err)))
				if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
					return
				}
				return
			}
			if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("authenticate gateway: %s", claim.Gateway)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
