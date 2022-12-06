package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"

	"github.com/microservice-sample-go/shared"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var input LoginPayload
	if body, err := io.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError("Invalid body: %s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	} else if err := json.Unmarshal(body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError("Invalid body: %s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(validationErrors.Error()))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", validationErrors.Error())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	userData, err := models.User.GetByEmail(input.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	if userData == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(shared.WriteError("%s", "email is not registered"))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: email is not registered;\nemail: %s", input.Email)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	passwordMatch, err := PasswordMatches(input.Password, userData.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	if !passwordMatch {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError("%s", "incorrect password"))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: incorrect password;\nemail: %s", input.Email)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	jwtToken, err := GenerateToken(GetConfig().JWTSecretKey, input.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	type Token struct {
		Token string `json:"token"`
	}
	if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("info: login successfully;\nemail: %v", input.Email)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		return
	}
	token := Token{Token: jwtToken}
	json.NewEncoder(w).Encode(token)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(shared.WriteInfo("%s says pong!", shared.AUTH_SERVICE))
}
