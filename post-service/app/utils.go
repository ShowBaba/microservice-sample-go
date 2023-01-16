package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/showbaba/microservice-sample-go/shared"
)

func ValidateGatewayToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("gateway_signature")
			if token == "" {
				shared.Dispatch400Error(w, "gateway-token not in header", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, "gateway-token not in header"); err != nil {
					shared.Dispatch500Error(w, err)
					return
				}
				return
			}
			claim, err := shared.ValidateGatewayToken(token, GetConfig().BlogServiceSecretKey)
			if err != nil {
				shared.Dispatch400Error(w, "error validating gateway-token", fmt.Sprintf("%v", err))
				if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
					shared.Dispatch500Error(w, err)
					return
				}
				return
			}
			if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("authenticate gateway: %s", claim.Gateway)); err != nil {
				shared.Dispatch500Error(w, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ValidateAuthToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				shared.Dispatch400Error(w, "auth token not in header", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, "auth token not in header"); err != nil {
					shared.Dispatch500Error(w, err)
					return
				}
				return
			}
			parts := strings.SplitN(authHeader, " ", 3)
			if len(parts) != 3 || parts[0] != "Bearer" {
				shared.Dispatch400Error(w, "bearer token not in header", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, "bearer token not in header"); err != nil {
					shared.Dispatch500Error(w, err)
					return
				}
				return
			}
			claim, err := shared.ValidateAuthToken(parts[2], GetConfig().JWTSecretKey)
			if err != nil {
				shared.Dispatch400Error(w, "error validating auth token token", fmt.Sprintf("%v", err))
				if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
					shared.Dispatch500Error(w, err)
					return
				}
				return
			}
			if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("authenticate user: %s", claim.Email)); err != nil {
				shared.Dispatch500Error(w, err)
				return
			}
			// set values in the request context
			ctx := context.WithValue(r.Context(), keyEmail, claim.Email)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

type key string

const (
	keyEmail key = "email"
)
