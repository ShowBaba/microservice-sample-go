package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/showbaba/microservice-sample-go/shared"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	requestInfo := ParseRequest(r)
	config, err := ReadConfig()
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	// get service name from url
	serviceName := strings.Split(requestInfo.OriginalURL, "/")[1]
	var requestParam string
	if len(strings.Split(requestInfo.OriginalURL, "/")) > 2 {
		requestParam = strings.Split(strings.Split(requestInfo.OriginalURL, "/")[2], "?")[0]
	}
	var request Request
	var service string
	switch serviceName {
	case "auth":
		// check path is available
		switch requestInfo.Method {
		case "GET":
			if len(config.Services.Auth.Endpoints.Get) < 1 {
				shared.Dispatch501Error(w, "request path is unavailable", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					shared.Dispatch400Error(w, "invalid body: %s", err)
					return
				}
				return
			} else {
				flag := false
				for _, param := range config.Services.Auth.Endpoints.Get {
					if param == requestParam {
						flag = true
					}
				}
				if !flag {
					shared.Dispatch501Error(w, "request path mismatch", nil)
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						shared.Dispatch400Error(w, "invalid body: %s", err)
						return
					}
					return
				}
				request.Method = "GET"
			}
		case "POST":
			if len(config.Services.Auth.Endpoints.Post) < 1 {
				shared.Dispatch501Error(w, "request path is unavailable", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					shared.Dispatch400Error(w, "invalid body: %s", err)
					return
				}
				return
			} else {
				flag := false
				for _, param := range config.Services.Auth.Endpoints.Post {
					if param == requestParam {
						flag = true
					}
				}
				if !flag {
					shared.Dispatch501Error(w, "request path mismatch", nil)
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						shared.Dispatch400Error(w, "invalid body: %s", err)
						return
					}
					return
				}
				request.Method = "POST"
			}
		default:
			shared.Dispatch405Error(w, "request path mismatch", nil)
			if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
				shared.Dispatch500Error(w, err)
				return
			}
			return
		}
		request.URL = fmt.Sprintf("%s:%v%s", config.Services.Auth.BaseURL, config.Services.Auth.Port, requestInfo.OriginalURL)
		service = shared.AUTH_SERVICE
	case "user":
		switch requestInfo.Method {
		case "GET":
			if len(config.Services.User.Endpoints.Get) < 1 {
				shared.Dispatch501Error(w, "request path is unavailable", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					shared.Dispatch400Error(w, "invalid body: %s", err)
					return
				}
				return
			} else {
				flag := false
				for _, param := range config.Services.User.Endpoints.Get {
					if param == requestParam {
						flag = true
					}
				}
				if !flag {
					shared.Dispatch501Error(w, "request path mismatch", nil)
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						shared.Dispatch400Error(w, "invalid body: %s", err)
						return
					}
					return
				}
				request.Method = "GET"
			}
		case "POST":
			if len(config.Services.User.Endpoints.Post) < 1 {
				shared.Dispatch501Error(w, "request path is unavailable", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					shared.Dispatch400Error(w, "invalid body: %s", err)
					return
				}
				return
			} else {
				flag := false
				for _, param := range config.Services.User.Endpoints.Post {
					if param == requestParam {
						flag = true
					}
				}
				if !flag {
					shared.Dispatch501Error(w, "request path mismatch", nil)
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						shared.Dispatch400Error(w, "invalid body: %s", err)
						return
					}
					return
				}
				request.Method = "POST"
			}
		default:
			shared.Dispatch405Error(w, "request method is not found", nil)
			if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request method is not found"); err != nil {
				shared.Dispatch500Error(w, err)
				return
			}
			return
		}
		request.URL = fmt.Sprintf("%s:%v%s", config.Services.User.BaseURL, config.Services.User.Port, requestInfo.OriginalURL)
		service = shared.USER_SERVICE
	case "post":
		switch requestInfo.Method {
		case "GET":
			if len(config.Services.Post.Endpoints.Get) < 1 {
				shared.Dispatch501Error(w, "request path is unavailable", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					shared.Dispatch400Error(w, "invalid body: %s", err)
					return
				}
				return
			} else {
				flag := false
				for _, param := range config.Services.Post.Endpoints.Get {
					if param == requestParam {
						flag = true
					}
				}
				if !flag {
					shared.Dispatch501Error(w, "request path mismatch", nil)
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						shared.Dispatch400Error(w, "invalid body: %s", err)
						return
					}
					return
				}
				request.Method = "GET"
			}
		case "POST":
			if len(config.Services.Post.Endpoints.Post) < 1 {
				shared.Dispatch501Error(w, "request path is unavailable", nil)
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					shared.Dispatch400Error(w, "invalid body: %s", err)
					return
				}
				return
			} else {
				flag := false
				for _, param := range config.Services.Post.Endpoints.Post {
					if param == requestParam {
						flag = true
					}
				}
				if !flag {
					shared.Dispatch501Error(w, "request path mismatch", nil)
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						shared.Dispatch400Error(w, "invalid body: %s", err)
						return
					}
					return
				}
				request.Method = "POST"
			}
		default:
			shared.Dispatch405Error(w, "request method is not found", nil)
			if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request method is not found"); err != nil {
				shared.Dispatch500Error(w, err)
				return
			}
			return
		}
		request.URL = fmt.Sprintf("%s:%v%s", config.Services.Post.BaseURL, config.Services.Post.Port, requestInfo.OriginalURL)
		service = shared.POST_SERVICE
	default:
		shared.Dispatch400Error(w, "no service available to process request", nil)
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: no service available to process request"); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	token, err := GenerateGatewayToken(service)
	if err != nil {
		shared.Dispatch500Error(w, err)   
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	request.GatewayAuth = token
	request.Auth = requestInfo.Auth
	body, err := ioutil.ReadAll(requestInfo.Body)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	response, err := ForwardRequest(&request, &body)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	w.WriteHeader(response.Status)
	responseJSON, err := json.Marshal(shared.APIResponse{
		Status:  (*response).Status,
		Message: (*response).Message,
		Data:    (*response).Data,
	})
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
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
		Message: fmt.Sprintf("%s says pong!", shared.GATEWAY_SERVICE),
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		shared.Dispatch500Error(w, err)
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			shared.Dispatch500Error(w, err)
			return
		}
		return
	}
	w.Write(responseJSON)
}
