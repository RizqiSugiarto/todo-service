package rabbitclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Config struct {
		Host           string `mapstructure:"host"`
		Port           int    `mapstructure:"port"`
		User           string `mapstructure:"user"`
		Pass           string `mapstructure:"pass"`
		Queues         string `mapstructure:"queues"`
		JobQueueSize   int    `mapstructure:"job_queue_size"`
		PublisherCount int    `mapstructure:"publisher_count"`
		ConsumerCount  int    `mapstructure:"consumer_count"`
	}

	RabbitMQ struct {
		Conn       *amqp.Connection
		Ch         *amqp.Channel
		jobs       chan job
		publishers chan struct{}
	}

	job struct {
		ctx     context.Context
		payload any
		queue   string
	}
)

func New(cfg Config) (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Pass, cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err.Error())
	}

	log.Printf("RabbitMQ connected on: %s", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	rabbitMQ := &RabbitMQ{
		Conn: conn,
	}

	return rabbitMQ, nil
}

func (r *RabbitMQ) NewChannel() error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %s", err.Error())
	}

	r.Ch = ch

	return nil
}

func (r *RabbitMQ) NewPublisherPool(cfg Config) {
	r.jobs = make(chan job, cfg.JobQueueSize)
	r.publishers = make(chan struct{}, cfg.PublisherCount)

	for i := 0; i < cfg.PublisherCount; i++ {
		r.publishers <- struct{}{}
	}

	go r.run()
}

func (r *RabbitMQ) DeclareQueue(queueName string) error {
	_, err := r.Ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare %s queue: %s", queueName, err.Error())
	}

	return nil
}

func (r *RabbitMQ) run() {
	for j := range r.jobs {
		<-r.publishers
		go func(j job) {
			defer func() { r.publishers <- struct{}{} }()
			err := r.publishMessage(j.ctx, j.payload, j.queue)
			if err != nil {
				log.Println("Failed to publish message:", err.Error())
			}
		}(j)
	}
}

func (r *RabbitMQ) publishMessage(ctx context.Context, payload any, name string) error {
	q, err := r.Ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %s", err.Error())
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %s", err.Error())
	}

	err = r.Ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %s", err.Error())
	}

	return nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, payload any, name string) {
	r.jobs <- job{
		ctx:     ctx,
		payload: payload,
		queue:   name,
	}
}

func (r *RabbitMQ) StartConsuming(queueName string, consumerCount int, handler func(amqp.Delivery)) error {
	msgs, err := r.Ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %s", err.Error())
	}

	for i := 0; i < consumerCount; i++ {
		go func() {
			for d := range msgs {
				handler(d)
			}
		}()
	}

	return nil
}
