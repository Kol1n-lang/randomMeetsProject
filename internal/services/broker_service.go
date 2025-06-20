package services

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"randomMeetsProject/config"
)

func NewPublisher(queue string, bodyMessage string) error {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		return err
	}
	conn, err := amqp.Dial(cfg.RabbitMQUrl())
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	que, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.Publish("", que.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(bodyMessage),
	})
	if err != nil {
		return err
	}
	return nil
}
