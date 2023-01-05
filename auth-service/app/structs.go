package app

import (
	"github.com/gorilla/mux"
)

type LoginPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type App struct {
	Router *mux.Router
}
