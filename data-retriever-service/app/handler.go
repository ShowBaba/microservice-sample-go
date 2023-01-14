package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/graphql-go/graphql"
)

type App struct {
	db *sql.DB
}

var (
	ctx = context.Background()
)

func NewApp(dbConn *sql.DB) *App {
	return &App{
		db: dbConn,
	}
}

type GraphQLPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (a *App) Init() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(UserType),
				Args: graphql.FieldConfigArgument{
					"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
					"offset": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var users []interface{}
					var result ListResult
					query := `SELECT * FROM Users`
					row := a.db.QueryRowContext(ctx, query)
					if err := row.Scan(users); err != nil {
						return result, err
					}
					fmt.Println("result - ", users)
					result.Nodes = users
					result.TotalCount = len(users)
					return result, nil
				},

				// "user": &graphql.Field{
				// 	Type: makeNodeListType("SingleUserList", UserType),
				// 	Args: graphql.FieldConfigArgument{
				// 		"id": &graphql.ArgumentConfig{
				// 			Type: graphql.Int,
				// 		},
				// 		"email": &graphql.ArgumentConfig{
				// 			Type: graphql.String,
				// 		},
				// 	},
				// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// 		var result ListResult
				// 		var user interface{}
				// 		id, idOk := params.Args["id"].(int)
				// 		email, emailOk := params.Args["email"].(string)
				// 		var row *sql.Row
				// 		if emailOk && idOk {
				// 			query := `SELECT * FROM users WHERE id = $1 AND email = $2`
				// 			row = a.db.QueryRowContext(ctx, query, id, email)
				// 		} else if idOk {
				// 			query := `SELECT * FROM users WHERE id = $1`
				// 			row = a.db.QueryRowContext(ctx, query, id)
				// 		} else if emailOk {
				// 			query := `SELECT * FROM users WHERE email = $1`
				// 			row = a.db.QueryRowContext(ctx, query, email)
				// 		}
				// 		if err := row.Scan(user); err != nil {
				// 			return result, err
				// 		}
				// 		result.Nodes = append(result.Nodes, user)
				// 		result.TotalCount = len(result.Nodes)
				// 		return result, nil
				// 	},
				// },
			},
		}})
}

// var RootQuery =

func makeListField(listType graphql.Output, resolve graphql.FieldResolveFn) *graphql.Field {
	listField := &graphql.Field{
		Type:    listType,
		Resolve: resolve,
		Args: graphql.FieldConfigArgument{
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
		},
	}
	return listField
}

func makeNodeListType(name string, nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"nodes":      &graphql.Field{Type: graphql.NewList(nodeType)},
			"totalCount": &graphql.Field{Type: graphql.Int},
		},
	})
}
