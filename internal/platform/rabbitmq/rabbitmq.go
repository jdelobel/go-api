package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ErrInvalidRabbitMQProvided is returned in the event that an uninitialized db is
// used to perform actions against.
var ErrInvalidRabbitMQProvided = errors.New("invalid RabbitMQ provided")

// RabbitMQ structure
type RabbitMQ struct {
	channel      *amqp.Channel
	defaultQueue *string
}

// NewRabbitMQ initialize a new RabbitMQ connection
func NewRabbitMQ(url string, defaultQueue *string) (*RabbitMQ, error) {
	rabbitCloseError := make(chan *amqp.Error)
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to dial rabbitmq: %v", err)
	}

	connection.NotifyClose(rabbitCloseError)

	channel, err := connection.Channel()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create rabbitmq channel: %v", err)
	}
	if defaultQueue != nil {
		_, err := channel.QueueDeclare(
			*defaultQueue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to create queue: %v", err)
		}
	}
	return &RabbitMQ{channel: channel, defaultQueue: defaultQueue}, nil
}

// DeclareQueue declare a queue
func (rbmq *RabbitMQ) DeclareQueue(queueName string) (*amqp.Queue, error) {
	if rbmq == nil {
		return nil, errors.Wrap(ErrInvalidRabbitMQProvided, "rbmq == nil")
	}
	queue, err := rbmq.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to create queue: %v", err)
	}
	return &queue, nil
}

// Publish send a jsonMessage to the queue
func (rbmq *RabbitMQ) Publish(queueName *string, jsonMessage []byte) error {
	if rbmq == nil {
		return errors.Wrap(ErrInvalidRabbitMQProvided, "rbmq == nil")
	}
	if queueName == nil && rbmq.defaultQueue == nil {
		return errors.New("Unable to send message because no queue is defined")
	} else if queueName == nil && rbmq.defaultQueue != nil {
		queueName = rbmq.defaultQueue
	}
	err := rbmq.channel.Publish(
		"",
		*queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMessage,
		})
	if err != nil {
		return errors.Wrapf(err, "Unable to send message: %v", err)
	}
	return nil
}
