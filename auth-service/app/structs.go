package app

import (
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type LoginPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

type App struct {
	Router *mux.Router
}
