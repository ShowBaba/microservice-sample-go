package data

import (
	"database/sql"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

type Models struct {
	User User
}

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}
