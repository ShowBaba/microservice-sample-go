package app

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/microservice-sample-go/shared"
	"golang.org/x/crypto/bcrypt"
)

func PasswordMatches(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// GenerateToken generates a jwt token
func GenerateToken(JWTSecretKey, email string) (signedToken string, err error) {
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return
	}
	return
}

func ValidateGatewayToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("gateway_signature")
			if token == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(shared.WriteError(http.StatusBadRequest, "gateway-token not in header", nil))
				if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, "gateway-token not in header"); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
					return
				}
				return
			}
			claim, err := shared.ValidateGatewayToken(token, GetConfig().AuthServiceSecretKey)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(shared.WriteError(http.StatusBadRequest, "error validating gateway-token", fmt.Sprintf("%v", err)))
				if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
					return
				}
				return
			}
			if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("authenticate gateway: %s", claim.Gateway)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
