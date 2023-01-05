package app

import "github.com/gorilla/mux"

type App struct {
	Router *mux.Router
}

type CreatePostPayload struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}
