package shared

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func WriteError(format string, args ...interface{}) []byte {
	response := map[string]string{
		"error": fmt.Sprintf(format, args...),
	}
	if data, err := json.Marshal(response); err == nil {
		return data
	} else {
		log.Printf("Err: %s", err)
	}
	return nil
}

func WriteInfo(format string, args ...interface{}) []byte {
	response := map[string]string{
		"info": fmt.Sprintf(format, args...),
	}
	if data, err := json.Marshal(response); err == nil {
		return data
	} else {
		log.Printf("Err: %s", err)
	}
	return nil
}

func ConnectToSQLDB(host, user, password, dbname string, port int) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("db successfully connected!")
	return db
}

func LogRequest(ctx context.Context, messageChan *amqp.Channel, source, data string) error {
	messageData := LogPayload{
		Source: source,
		Data:        data,
	}
	b, err := json.Marshal(messageData)
	if err != nil {
		log.Fatal(err)
	}
	// create a message to publish
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(b),
	}
	// publish message to queue
	if err = messageChan.PublishWithContext(
		ctx,
		"",
		LOGGER_SERVICE,
		false,
		false,
		message,
	); err != nil {
		return err
	}
	return nil
}
