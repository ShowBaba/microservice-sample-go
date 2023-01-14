package app

import "github.com/graphql-go/graphql"

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"Id":         &graphql.Field{Type: graphql.String},
			"email":      &graphql.Field{Type: graphql.String},
			"firstname":  &graphql.Field{Type: graphql.String},
			"lastname":   &graphql.Field{Type: graphql.String},
			"created_at": &graphql.Field{Type: graphql.DateTime},
			"updated_at": &graphql.Field{Type: graphql.DateTime},
		},
	},
)

var PostType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"Id":         &graphql.Field{Type: graphql.String},
			"title":      &graphql.Field{Type: graphql.String},
			"body":       &graphql.Field{Type: graphql.String},
			"user_id":    &graphql.Field{Type: graphql.Int},
			"user":       &graphql.Field{Type: UserType},
			"created_at": &graphql.Field{Type: graphql.DateTime},
			"updated_at": &graphql.Field{Type: graphql.DateTime},
		},
	},
)
