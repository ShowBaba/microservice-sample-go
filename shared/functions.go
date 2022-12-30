package shared

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func WriteError(statusCode int, message string, err interface{}) []byte {
	response := APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    err,
	}
	data, err := json.Marshal(response)
	if err == nil {
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
		Data:   data,
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

func ValidateGatewayToken(signedToken, SECRET_KEY string) (*GatewayTokenJwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&GatewayTokenJwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*GatewayTokenJwtClaim)
	if !ok {
		return nil, err
	}
	// check the expiration date of the token
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, err
	}
	return claims, nil
}

// publish to notifications topic
func SendNotification(channel *amqp.Channel, payload []byte) error {
	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = channel.QueueBind(q.Name, q.Name, NOTIFICATION_TOPIC, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Set up a channel to receive acknowledgement messages
	// ackMsgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	err = channel.Publish(NOTIFICATION_TOPIC, "email.welcome", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        payload,
		ReplyTo:     q.Name,
	})
	if err != nil {
		log.Fatal(err)
	}
	// for {
	// 	// Wait for an acknowledgement message
	// 	select {
	// 	case <-ackMsgs:
	// 		log.Println("Received acknowledgement from notification subscriber service")
	// 		return nil
	// 	default:
	// 		log.Println("waiting for ack")
	// 		continue
	// 	}
	// }

	log.Println("Sent notifications")
	return nil
}
