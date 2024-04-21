package model

import (
	"context"
	"github.com/go-faker/faker/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func Produce(connectionString string, exchange string, message string, key string) {
	conn, err := amqp.Dial("amqp://" + connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	err = ch.Confirm(false)
	if err != nil {
		log.Fatalf("channel.confirm: %s", err)
	}
	defer ch.Close()

	// Confirms the message has been received
	confirms := make(chan amqp.Confirmation, 1)
	confirms = ch.NotifyPublish(confirms)

	// Confirms the message has been received and routed
	returns := make(chan amqp.Return)
	returns = ch.NotifyReturn(returns)

	// See the Channel.Consume example for the complimentary declare.
	err = ch.ExchangeDeclarePassive(exchange, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange.declare: %v", err)
	}

	msg := amqp.Publishing{
		Timestamp:   time.Now(),
		ContentType: "text/plain",
	}

	if message != "" {
		msg.Body = []byte(message)
	} else {
		msg.Body = []byte(faker.Sentence())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx, exchange, key, true, false, msg)
	if err != nil {
		// Since publish is asynchronous this can happen if the network connection
		// is reset or if the server has run out of resources.
		log.Fatalf("basic.publish: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel2()

	select {
	case ret := <-returns:
		log.Println("Message: ", string(msg.Body), " NOT sent or NOT routed returned code: ", ret.ReplyCode, " Reply text: ", ret.ReplyText)
		cancel2()
	case confirm := <-confirms:
		if confirm.Ack {
			log.Println("Acknowledged message: ", string(msg.Body))
		} else {
			log.Println("message: ", string(msg.Body), " NOT routed")
		}

	case <-ctx2.Done():
		log.Println("Message: ", string(msg.Body), " NOT sent or NOT routed timeout")
	}

	err = ch.Close()
	if err != nil {
		log.Println("Failed to close channel err: ", err)
	}
}
