package service

import amqp "github.com/rabbitmq/amqp091-go"

type rabbitConfig struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func dialRabbit() (*rabbitConfig, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"user",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &rabbitConfig{conn: conn, ch: ch, q: q}, nil
}
