package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"

	"github.com/microservice-sample-go/shared"
	"github.com/microservice-sample-go/user-service/data"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var input RegisterPayload
	if body, err := io.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(http.StatusBadRequest, "Invalid body: %s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	} else if err := json.Unmarshal(body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(http.StatusBadRequest, "Invalid body: %s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", validationErrors.Error())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	userData, err := models.User.GetByEmail(input.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	if userData != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(http.StatusBadRequest, "email already used", nil))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: email already used;\nemail: %s", input.Email)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	hash, err := HashPassword(input.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	user := data.User{
		Email:     input.Email,
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Password:  hash,
	}
	id, err := models.User.Insert(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	// TODO: send id to avater-generator function
	// TODO: trigger notification service, send email notification to user email
	if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("info: register successfully;\nemail: %v", input.Email)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		return
	}
	response := shared.APIResponse{
		Status:  http.StatusOK,
		Message: "user registered",
		Data:    map[string]string{"id": fmt.Sprint(id)},
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
		Message: fmt.Sprintf("%s says pong!", shared.USER_SERVICE),
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	w.Write(responseJSON)
}
