package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func connectRabbitMQ() *amqp091.Connection {
	for i := 0; i < 30; i++ {
		conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			log.Println("RabbitMQ connected")
			return conn
		}
		log.Printf("RabbitMQ wait %ds: %v", i*2, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("RabbitMQ unavailable")
	return nil
}

func main() {
	conn := connectRabbitMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("bugs", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Queue ready: %s (%d messages)", q.Name, q.Messages)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintf(w, "Bug service: queue=%d", q.Messages)
		})
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Println("Bug service consuming bugs...")
	for msg := range msgs {
		var bug map[string]interface{}
		json.Unmarshal(msg.Body, &bug)
		log.Printf("✅ Processed bug: %v", bug["description"])
	}
}
