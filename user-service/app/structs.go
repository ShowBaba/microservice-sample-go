package app

import "github.com/gorilla/mux"

type RegisterPayload struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type App struct {
	Router *mux.Router
}
