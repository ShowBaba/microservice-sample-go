package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator"

	"github.com/showbaba/microservice-sample-go/shared"
	"github.com/showbaba/microservice-sample-go/user-service/data"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var input RegisterPayload
	if body, err := io.ReadAll(r.Body); err != nil {
		shared.Dispatch400Error(w, "invalid body: %s", err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	} else if err := json.Unmarshal(body, &input); err != nil {
		shared.Dispatch400Error(w, "invalid body: %s", err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
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
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", validationErrors.Error())); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	userData, err := models.User.GetByEmail(input.Email)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	if userData != nil {
		shared.Dispatch400Error(w, "email already used", nil)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: email already used;\nemail: %s", input.Email)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	hash, err := HashPassword(input.Password)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
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
	id, err := user.Insert()
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	// TODO: send id to avater-generator function
	mail := shared.Mail{
		Sender:  shared.MAIL_USERNAME,
		Subject: "Welcome to our blog!",
		To:      []string{input.Email},
		Body: `<div style="font-family: Helvetica, Arial, sans-serif; min-width: 1000px; overflow: auto; line-height: 2;">
            <div style="margin: 50px auto; width: 70%; padding: 20px 0;">
                <div style="border-bottom: 1px solid #eee;"><a href="blog.com" style="font-size: 1.4em; color: #00466a; text-decoration: none; font-weight: 600;">SAM's BLOG</a></div>
                <p style="font-size: 1.1em;">Hi,</p>
                <p>Hi ` + input.Firstname + `</p>
                <p>Welcome to Sam's BLOG</p>
                <p style="font-size: 0.9em;">
                    Regards,<br />
                    SAM's BLOG
                </p>
                <hr style="border: none; border-top: 1px solid #eee;" />
            </div>
        </div>`,
	}
	payload, err := json.Marshal(mail)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	if err := shared.SendNotification(messageChan, payload); err != nil {
		if err != nil {
			shared.Dispatch500Error(w, err)
			if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
				shared.Dispatch500Error(w, err)
				return
			}
			return
		}
	}
	response := shared.APIResponse{
		Status:  http.StatusOK,
		Message: "user registered",
		Data:    map[string]string{"id": fmt.Sprint(id)},
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
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
		Message: fmt.Sprintf("%s says pong!", shared.USER_SERVICE),
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.USER_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	w.Write(responseJSON)
}
