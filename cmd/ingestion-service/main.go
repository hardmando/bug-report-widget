package main

import (
	"fmt"
	"io"
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
		log.Printf("RabbitMQ wait %d/30s: %v", i*2, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("RabbitMQ timeout")
	return nil
}

func main() {
	conn := connectRabbitMQ()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Channel:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("bugs", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Queue:", err)
	}
	log.Printf("Queue: %s (%d msgs)", q.Name, q.Messages)

	http.HandleFunc("/ingest/bugs", func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-API-Key, Content-Type")
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Ingested bug: %s", body)

		err = ch.Publish("", q.Name, false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		if err != nil {
			log.Printf("Publish error: %v", err)
			http.Error(w, "Queue error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{"status":"queued"}`)
	})

	log.Println("Ingestion service on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
