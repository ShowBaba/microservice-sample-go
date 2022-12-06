package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbCLient *mongo.Client

type Models struct {
	Log Log
}

func NewMongoClient(client *mongo.Client) Models {
	dbCLient = client
	return Models{
		Log: Log{},
	}
}

func ConnectToDB(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, ctx, cancel, err
	}
	return client, ctx, cancel, err
}

func PingDB(ctx context.Context, client *mongo.Client) {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("db connected successfully")
}

func CloseDBConnection(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
