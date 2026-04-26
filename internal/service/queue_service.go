package service

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewQueueService(url string) (*QueueService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &QueueService{conn: conn, ch: ch}, nil
}

func (q *QueueService) Publish(queue string, body []byte) error {
	if _, err := q.ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}
	return q.ch.PublishWithContext(context.Background(), "", queue, false, false,
		amqp.Publishing{ContentType: "application/json", Body: body},
	)
}

func (q *QueueService) Consume(queue string, handler func([]byte)) error {
	if _, err := q.ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}
	msgs, err := q.ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()
	return nil
}

func (q *QueueService) Close() {
	q.ch.Close()
	q.conn.Close()
}
