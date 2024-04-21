package model

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Consume(connectionString string, exchange string, queue string, key string, durable bool) {
	conn, err := amqp.Dial("amqp://" + connectionString)
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
	err = ch.ExchangeDeclarePassive(exchange, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange.declare: %v", err)
	}

	// if the queue is transient then we auto-delete it
	if !durable {
		_, err = ch.QueueDeclare(queue, false, true, false, false, nil)
		if err != nil {
			log.Fatalf("queue.declare: %v", err)
		}
	} else {
		_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			log.Fatalf("queue.declare: %v", err)
		}
	}

	err = ch.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		log.Fatalf("queue.bind: %v", err)
	}

	// Set our quality of service.  Since we're sharing 1 consumer1 on the same
	// channel, we want at least 1 messages in flight.
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("basic.qos: %v", err)
	}

	msg, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case m := <-msg:
				if e := m.Ack(false); e != nil {
					log.Printf("ack error: %+v", e)
				}
				log.Printf("Message: %s \n", string(m.Body))

			case <-sigs:
				log.Printf("Received SIGQUIT, shutting down")
				goto end
			}
		}
	end:
		wg.Done()
	}()

	wg.Wait()

	if !durable {
		_, err = ch.QueueDelete(queue, false, false, true)
		if err != nil {
			log.Printf("Problem occurred during deletion: %v", err)
		}
	}

	ch.Close()

	log.Println("Exiting Consumer")
}
