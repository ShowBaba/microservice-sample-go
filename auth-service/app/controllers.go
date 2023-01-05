package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"

	"github.com/showbaba/microservice-sample-go/shared"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var input LoginPayload
	if body, err := io.ReadAll(r.Body); err != nil {
		shared.Dispatch400Error(w, "invalid body: %s", err)
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	} else if err := json.Unmarshal(body, &input); err != nil {
		shared.Dispatch400Error(w, "invalid body: %s", err)
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(http.StatusBadRequest, validationErrors.Error(), nil))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", validationErrors.Error())); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	userData, err := models.User.GetByEmail(input.Email)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)

			return
		}
		return
	}
	if userData == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(shared.WriteError(http.StatusNotFound, "email is not registered", nil))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: email is not registered;\nemail: %s", input.Email)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	passwordMatch, err := PasswordMatches(input.Password, userData.Password)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	if !passwordMatch {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(http.StatusBadRequest, "incorrect password", nil))
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: incorrect password;\nemail: %s", input.Email)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	jwtToken, err := GenerateToken(GetConfig().JWTSecretKey, input.Email)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	type Token struct {
		Token string `json:"token"`
	}
	if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("info: login successfully;\nemail: %v", input.Email)); err != nil {
		shared.Dispatch500Error(w, err)
		return
	}
	token := Token{Token: jwtToken}
	response := shared.APIResponse{
		Status:  http.StatusOK,
		Message: "login",
		Data:    token,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	w.Write(responseJSON)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	response := shared.APIResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("%s says pong!", shared.AUTH_SERVICE),
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.AUTH_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	w.Write(responseJSON)
}
