package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
	"log"
	"strings"
)

type EmailConfig struct {
	Sender      string
	AppPassword string
}

func MessageHandler(cfg EmailConfig) func(string, amqp.Delivery) error {
	return func(queueName string, msg amqp.Delivery) error {
		switch queueName {
		case "email-confirm":
			return sendEmail(cfg, msg)
		default:
			log.Printf("Unknown queue: %s", queueName)
			return nil
		}
	}
}

func sendEmail(cfg EmailConfig, msg amqp.Delivery) error {
	fmt.Println("We are register a new user")
	dialer := gomail.NewDialer("smtp.gmail.com", 587, cfg.Sender, cfg.AppPassword)
	email := gomail.NewMessage()
	toMail := strings.Split(string(msg.Body), " ")
	email.SetHeader("From", cfg.Sender)
	email.SetHeader("To", toMail[1])
	email.SetHeader("Subject", "MeetApp")
	email.SetBody("text/html", "127.0.0.1:8080/api/v1/auth/confirm-email/?user_id="+toMail[0])

	return dialer.DialAndSend(email)
}
