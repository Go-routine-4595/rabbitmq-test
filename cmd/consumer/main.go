package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@34.172.210.10:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// See the Channel.Consume example for the complimentary declare.
	err = ch.ExchangeDeclarePassive("hw", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange.declare: %v", err)
	}

	_, err = ch.QueueDeclare("alarms", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue.declare: %v", err)
	}
	err = ch.QueueBind("alarms", "", "hw", false, nil)
	if err != nil {
		log.Fatalf("queue.bind: %v", err)
	}

	// Set our quality of service.  Since we're sharing 1 consumer1 on the same
	// channel, we want at least 1 messages in flight.
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("basic.qos: %v", err)
	}

	msg, err := ch.Consume("alarms", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for m := range msg {
			// ... this consumer is responsible for sending emails per log
			if e := m.Ack(false); e != nil {
				log.Printf("ack error: %+v", e)
			}
			log.Printf("Message: %s \n", string(m.Body))
		}
		wg.Done()
	}()
	wg.Wait()
	log.Println("Exiting Consumer")
}
