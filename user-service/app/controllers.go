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
		w.Write(shared.WriteError("Invalid body: %s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	} else if err := json.Unmarshal(body, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError("Invalid body: %s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
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
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", validationErrors.Error())); err != nil {
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
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	if userData != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError("email already used"))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: email already used;\nemail: %s", input.Email)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	hash, err := HashPassword(input.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
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
		w.Write(shared.WriteError("%s", err))
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError("%s", err))
			return
		}
		return
	}
	// TODO: send id to avater-generator function
	// TODO: trigger notification service, send email notification to user email
	if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("info: register successfully;\nemail: %v", input.Email)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError("%s", err))
		return
	}
	w.Write(shared.WriteInfo("user registered successfully with id: %s", id))
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(shared.WriteInfo("%s says pong!", shared.USER_SERVICE))
}
