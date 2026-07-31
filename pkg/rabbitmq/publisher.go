package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Publisher(ctx context.Context, body any, routingKey string) error {
	if r.ch.IsClosed() || r.conn.IsClosed(){
		err := r.connect()
		if err != nil {
			return fmt.Errorf("r.connect: %w",err)
		}
	}

	bytesBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w",err)
	}

	if err = r.ch.PublishWithContext(
		ctx,
		"",
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body: bytesBody,
		},
	); err != nil {
		return fmt.Errorf("r.Ch.PublishWithContext: %w",err)
	}
	return nil
}