package queue

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type rabbitMq struct {
	clientPath   string
	exchangeName string
	queueName    string

	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

type RabbitMQ interface {
	InitQueue(ctx context.Context) error
	Publish(ctx context.Context, dados []byte) error
	Consume(ctx context.Context) (<-chan amqp.Delivery, error)
	Close()
}

func NewRabbitMQ(clientPath string,
	exchangeName string,
	queueName string) RabbitMQ {

	return &rabbitMq{
		clientPath:   clientPath,
		exchangeName: exchangeName,
		queueName:    queueName}
}

func (r *rabbitMq) InitQueue(ctx context.Context) error {
	var err error

	fmt.Printf("INIT QUEUE %s\n", r.queueName)

	r.conn, err = amqp.Dial(r.clientPath)
	if err != nil {
		fmt.Printf("failed to connect to rabbitmq. %s\n", err.Error())
		return err
	}

	r.ch, err = r.conn.Channel()
	if err != nil {
		fmt.Printf("failed to create a channel rabbitmq. %s\n", err.Error())
		return err
	}

	err = r.ch.ExchangeDeclare(r.exchangeName, // name
		"direct", // type
		false,    // durable
		false,    // autodelete
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		fmt.Printf("failed to declare a exchange. %s\n", err.Error())
		return err
	}

	r.queue, err = r.ch.QueueDeclare(r.queueName, // name
		true,  // durable
		false, // autodelete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		fmt.Printf("failed to declare a queue. %s\n", err.Error())
		return err
	}

	err = r.ch.QueueBind(r.queue.Name, // queue name
		"",             // routing key
		r.exchangeName, // exchange
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *rabbitMq) Publish(ctx context.Context, dados []byte) error {

	err := r.ch.Publish(
		"",          // exchange
		r.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        dados, //[]byte(body)
		})

	if err != nil {
		fmt.Printf("failed to publish a message. %s\n", err.Error())
		return err
	}
	return nil
}

func (r *rabbitMq) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	var msgs <-chan amqp.Delivery
	msgs, err := r.ch.Consume(
		r.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		return msgs, err
	}
	return msgs, nil
}

func (r *rabbitMq) Close() {
	r.ch.Close()
	r.conn.Close()
}
