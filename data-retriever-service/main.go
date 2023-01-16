package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/graphql-go/graphql"
	app_ "github.com/showbaba/microservice-sample-go/data-retriever-service/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	schema graphql.Schema
	ctx    = context.Background()
)

func main() {
	var (
		err     error
		port, _ = strconv.ParseUint(strconv.Itoa(app_.GetConfig().DbPort), 10, 32)
		dsn     = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			app_.GetConfig().DbHost, port, app_.GetConfig().DbUser, app_.GetConfig().DbPassword, app_.GetConfig().DbName,
		)
	)
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println(
			err.Error(),
		)
		panic("failed to connect database")
	}
	app := app_.NewApp(dbConn)
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
	serverPort := app_.GetConfig().Port
	fmt.Println("GRAPHQL server is running on port ", serverPort)
	http.ListenAndServe(serverPort, nil)
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
		resp = graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: string(body),
			Context:       ctx,
		})
	}
	if len(resp.Errors) > 0 {
		responseError(w, fmt.Sprintf("%+v", resp.Errors), http.StatusBadRequest)
		return
	}
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
