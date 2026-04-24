package queue

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to Rabbitmq: %w", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open channel: %w", err)
	}
	return conn, channel, nil
}

func Publish(ch *amqp.Channel, queueName string, message []byte) error {
	_, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		"", queueName, false, false,
		amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func Consume(ch *amqp.Channel,queueName string) (<-chan amqp.Delivery,error){
	//ensure the queue exists
	_,err:=ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,
	)
	if err != nil {
		return nil,fmt.Errorf("failed to declare queue: %w", err)
	}
	
	//consume
	msgs,err:=ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	) 
	if err!=nil{
		return nil,fmt.Errorf("failed to register a consumer: %w", err)
	}
	return msgs,nil
}