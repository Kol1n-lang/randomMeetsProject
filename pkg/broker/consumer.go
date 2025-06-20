package broker

import (
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  []string
	handler func(queueName string, msg amqp.Delivery) error
}

func NewConsumer(amqpURL string, queues []string, handler func(string, amqp.Delivery) error) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queues:  queues,
		handler: handler,
	}, nil
}

func (c *Consumer) Start() {
	var wg sync.WaitGroup

	for _, queueName := range c.queues {
		wg.Add(1)

		go func(q string) {
			defer wg.Done()

			queue, err := c.channel.QueueDeclare(
				q,
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Printf("Failed to declare queue %s: %v", q, err)
				return
			}

			messages, err := c.channel.Consume(
				queue.Name,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				log.Printf("Failed to consume from %s: %v", q, err)
				return
			}

			log.Printf("Listening to queue: %s", q)
			for msg := range messages {
				if err := c.handler(q, msg); err != nil {
					log.Printf("Failed to process message from %s: %v", q, err)
					continue
				}
				msg.Ack(false)
			}
		}(queueName)
	}

	log.Println("All consumers started. Waiting for messages...")
	wg.Wait()
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}
