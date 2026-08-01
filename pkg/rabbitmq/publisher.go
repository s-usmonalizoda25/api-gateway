package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body any) error {
	if err := r.ensureConnected(); err != nil {
		return fmt.Errorf("ensureConnected: %w", err)
	}

	bytesBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	r.mu.Lock()
	ch := r.ch
	r.mu.Unlock()

	if err := ch.PublishWithContext(
		ctx,
		"",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         bytesBody,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		return fmt.Errorf("PublishWithContext: %w", err)
	}

	return nil
}
