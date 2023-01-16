package app

import (
	"github.com/graphql-go/graphql"
	"github.com/showbaba/microservice-sample-go/shared"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB
}

func NewApp(dbConn *gorm.DB) *App {
	return &App{
		db: dbConn,
	}
}

func (a *App) Init() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"users": makeListField(
				makeNodeListType("UserList", UserType),
				func(params graphql.ResolveParams) (interface{}, error) {
					var (
						list  ListResult
						tx    *gorm.DB
						users []shared.User
					)
					tx = parseDbClause(params, a.db.Model(&shared.User{}), UserType)
					res := tx.Debug().Scan(&users)
					if res.RowsAffected > 0 {
						list.Nodes = []interface{}{}
						for _, u := range users {
							list.Nodes = append(list.Nodes, interface{}(u))
						}
						list.TotalCount = len(list.Nodes)
					}
					return list, nil
				},
			),
			"posts": makeListField(
				makeNodeListType("PostList", PostType),
				func(params graphql.ResolveParams) (interface{}, error) {
					var (
						list  ListResult
						tx    *gorm.DB
						posts []shared.Post
					)
					tx = a.db.Model(&shared.Post{}).Joins("JOIN users ON posts.user_id = users.id")
					tx = parseDbClause(params, tx, PostType)
					res := tx.Debug().Scan(&posts)
					if res.RowsAffected > 0 {
						list.Nodes = []interface{}{}
						for _, u := range posts {
							list.Nodes = append(list.Nodes, interface{}(u))
						}
						list.TotalCount = len(list.Nodes)
					}
					return list, nil
				},
			),
		}})
}
