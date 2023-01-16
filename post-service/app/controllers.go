package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/showbaba/microservice-sample-go/post-service/data"
	"github.com/showbaba/microservice-sample-go/shared"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var input CreatePostPayload
	if body, err := io.ReadAll(r.Body); err != nil {
		shared.Dispatch400Error(w, "invalid body: %s", err)
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	} else if err := json.Unmarshal(body, &input); err != nil {
		shared.Dispatch400Error(w, "invalid body: %s", err)
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		shared.Dispatch400Error(w, "validation error", validationErrors.Error())
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", validationErrors.Error())); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	// get email from the request context
	email := r.Context().Value(keyEmail)
	userData, err := models.User.GetByEmail(email.(string))
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	if userData == nil {
		shared.Dispatch404Error(w, "user not found", fmt.Sprintf(`user with email (%s) not found`, email))
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf(`user with email (%s) not found`, email)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	post := data.Post{
		Title:  input.Title,
		Body:   input.Body,
		UserID: int64(userData.ID),
		User: shared.User{
			ID:    userData.ID,
			Email: userData.Email,
		},
	}
	id, err := post.Insert()
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	// TODO: send in app notification to user
	response := shared.APIResponse{
		Status:  http.StatusOK,
		Message: "post created",
		Data:    id,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
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
		Message: fmt.Sprintf("%s says pong!", shared.POST_SERVICE),
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.POST_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	w.Write(responseJSON)
}
