package rabbitmq

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const BookingQueue = "booking_queue"

type RabbitMQ struct {
	mu   sync.Mutex
	conn *amqp.Connection
	ch   *amqp.Channel
	url  string
}

func New(url string) (*RabbitMQ, error) {
	r := &RabbitMQ{url: url}

	if err := r.connect(); err != nil {
		return nil, fmt.Errorf("rabbitmq.New: %w", err)
	}

	return r, nil
}

func (r *RabbitMQ) connect() error {
	r.closeLocked()

	conn, err := amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("conn.Channel: %w", err)
	}

	if _, err := ch.QueueDeclare(
		BookingQueue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("ch.QueueDeclare: %w", err)
	}

	r.conn = conn
	r.ch = ch
	return nil
}

func (r *RabbitMQ) ensureConnected() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil && !r.conn.IsClosed() && r.ch != nil && !r.ch.IsClosed() {
		return nil
	}
	return r.connect()
}

func (r *RabbitMQ) closeLocked() {
	if r.ch != nil && !r.ch.IsClosed() {
		r.ch.Close()
	}
	if r.conn != nil && !r.conn.IsClosed() {
		r.conn.Close()
	}
}

func (r *RabbitMQ) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closeLocked()
}
