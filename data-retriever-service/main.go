package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
	app_ "github.com/showbaba/microservice-sample-go/data-retriever-service/app"
	"github.com/showbaba/microservice-sample-go/shared"
)

var (
	schema graphql.Schema
	ctx    = context.Background()
)

func main() {
	dbConn := shared.ConnectToSQLDB(
		app_.GetConfig().DbHost,
		app_.GetConfig().DbUser,
		app_.GetConfig().DbPassword,
		app_.GetConfig().DbName,
		app_.GetConfig().DbPort,
	)
	defer dbConn.Close()
	app := app_.NewApp(dbConn)
	var err error
	schema, err = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: app.Init(),
		},
	)
	if err != nil {
		fmt.Println("error creating schema: ", err)
		return
	}

	http.HandleFunc("/", Run)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func Run(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Read the query
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var (
		payload app_.GraphQLPayload
		resp    *graphql.Result
	)

	if err := json.Unmarshal(body, &payload); err == nil {
		// Perform GraphQL request
		resp = graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  payload.Query,
			VariableValues: payload.Variables,
			Context:        ctx,
		})
	} else {
		// Perform GraphQL request
		resp = graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: string(body),
			Context:       ctx,
		})
	}
	// Check for errors
	if len(resp.Errors) > 0 {
		responseError(w, fmt.Sprintf("%+v", resp.Errors), http.StatusBadRequest)
		return
	}
	// Return the result
	responseJSON(w, resp)
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
