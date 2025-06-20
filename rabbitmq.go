package main

import (
	"log"
	"randomMeetsProject/config"
	"randomMeetsProject/pkg/broker"
)

func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	queues := []string{"email-confirm"}

	emailCfg := broker.EmailConfig{
		Sender:      cfg.External.Sender,
		AppPassword: cfg.External.AppPassword,
	}

	consumer, err := broker.NewConsumer(
		cfg.RabbitMQUrl(),
		queues,
		broker.MessageHandler(emailCfg),
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	consumer.Start()
}
