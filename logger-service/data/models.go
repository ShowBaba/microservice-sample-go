package data

import (
	"context"
	"log"
	"time"
)

type Log struct {
	Data      string    `bson:"data" json:"data"`
	Source    string    `bson:"source" json:"source"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *Log) Insert(entry Log) error {
	collection := dbCLient.Database("blog-db").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), Log{
		Data:      entry.Data,
		Source:    entry.Source,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting into Logs", err)
		return err
	}
	return nil
}
