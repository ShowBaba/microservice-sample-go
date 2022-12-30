package app

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/microservice-sample-go/shared"
)

func HandleEmailMsg(ctx context.Context, channel *amqp.Channel, payload EmailMsgPayload) error {
	mail := shared.Mail{
		Sender:  shared.MAIL_USERNAME,
		Subject: payload.Subject,
		To:      payload.To,
		Body:    payload.Body,
	}
	log.Println("sending email to - ", mail.To)
	// if err := SendEmail(mail); err != nil {
	// 	if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("error while sending mail; err: %v", err)); err != nil {
	// 		return err
	// 	}
	// 	return fmt.Errorf("error while sending mail; err: %v", err)
	// }
	if err := shared.LogRequest(ctx, channel, shared.NOTIFICATION_SERVICE, fmt.Sprintf("sent email; payload: %v", payload)); err != nil {
		return err
	}
	return nil
}
