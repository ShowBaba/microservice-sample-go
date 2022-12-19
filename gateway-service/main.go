package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	app_ "github.com/microservice-sample-go/gateway-service/app"
	"github.com/microservice-sample-go/shared"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	messageChan *amqp.Channel
	ctx         = context.Background()
)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	// a.Router.Use(handler())
	a.Any("/{[^/]+}", handler)
	a.Any("/{[^/]+}/{[^/]+}", handler)
	a.Any("/{[^/]+}/{[^/]+}/{[^/]+}", handler)
	a.Get("/gateway/ping", Ping)
}

func (a *App) Any(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Get", "Post", "Patch", "Put")
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("Get")
}

func (a *App) Run(port string) {
	log.Fatal(
		http.ListenAndServe(
			port,
			handlers.CORS(
				handlers.AllowCredentials(),
				handlers.AllowedMethods([]string{"POST"}),
				handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
				handlers.MaxAge(3600),
			)(a.Router),
		),
	)
}

func main() {
	// setup rabbitmq
	connection, err := amqp.Dial(shared.RABBITMQ_SERVER_URL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	messageChan, err = connection.Channel()
	if err != nil {
		panic(err)
	}
	defer messageChan.Close()
	app := App{}
	port := app_.GetConfig().Port
	app.Initialize()
	log.Printf("starting server on port: %s", port)
	app.Run(port)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	requestInfo := app_.ParseRequest(r)
	config, err := app_.ReadConfig()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	// get service name from url
	serviceName := strings.Split(requestInfo.OriginalURL, "/")[1]
	requestParam := strings.Split(strings.Split(requestInfo.OriginalURL, "/")[2], "?")[0]
	var request app_.Request
	var service string
	switch serviceName {
	case "auth":
		// check path is available
		switch requestInfo.Method {
		case "GET":
			if len(config.Services.Auth.Endpoints.Get) < 1 {
				w.WriteHeader(http.StatusNotImplemented)
				w.Write(shared.WriteError(http.StatusNotImplemented, "request path is unavailable", nil))
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
					w.WriteHeader(http.StatusNotImplemented)
					w.Write(shared.WriteError(http.StatusNotImplemented, "request path mismatch", nil))
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
						return
					}
					return
				}
				request.Method = "GET"
			}
		case "POST":
			if len(config.Services.Auth.Endpoints.Post) < 1 {
				w.WriteHeader(http.StatusNotImplemented)
				w.Write(shared.WriteError(http.StatusNotImplemented, "request path is unavailable", nil))
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
					w.WriteHeader(http.StatusNotImplemented)
					w.Write(shared.WriteError(http.StatusNotImplemented, "request path mismatch", nil))
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
						return
					}
					return
				}
				request.Method = "POST"
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(shared.WriteError(http.StatusMethodNotAllowed, "request path mismatch", nil))
			if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
				w.WriteHeader(http.StatusNotImplemented)
				w.Write(shared.WriteError(http.StatusNotImplemented, "request path is unavailable", nil))
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
					w.WriteHeader(http.StatusNotImplemented)
					w.Write(shared.WriteError(http.StatusNotImplemented, "request path mismatch", nil))
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
						return
					}
					return
				}
				request.Method = "GET"
			}
		case "POST":
			if len(config.Services.User.Endpoints.Post) < 1 {
				w.WriteHeader(http.StatusNotImplemented)
				w.Write(shared.WriteError(http.StatusNotImplemented, "request path is unavailable", nil))
				if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path is unavailable"); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
					w.WriteHeader(http.StatusNotImplemented)
					w.Write(shared.WriteError(http.StatusNotImplemented, "request path mismatch", nil))
					if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request path mismatch"); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
						return
					}
					return
				}
				request.Method = "POST"
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(shared.WriteError(http.StatusMethodNotAllowed, "request method is not found", nil))
			if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: request method is not found"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
				return
			}
			return
		}
		request.URL = fmt.Sprintf("%s:%v%s", config.Services.User.BaseURL, config.Services.User.Port, requestInfo.OriginalURL)
		service = shared.USER_SERVICE
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(shared.WriteError(http.StatusBadRequest, "no service available to process request", nil))
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, "err: no service available to process request"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	token, err := app_.GenerateGatewayToken(service)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	request.GatewayAuth = token
	request.Auth = requestInfo.Auth
	body, err := ioutil.ReadAll(requestInfo.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	response, err := app_.ForwardRequest(&request, &body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
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
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
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
		Message: fmt.Sprintf("%s says pong!", shared.GATEWAY_SERVICE),
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
		if err := shared.LogRequest(ctx, messageChan, shared.GATEWAY_SERVICE, fmt.Sprintf("err: %v", err)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(shared.WriteError(http.StatusInternalServerError, "", fmt.Sprintf("%v", err)))
			return
		}
		return
	}
	w.Write(responseJSON)
}
