package main

import (
	"log"

	"github.com/streadway/amqp"
)

func newTask(jobID uint, request []byte, queueName string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	// body, err := json.Marshal(gradingRequest)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(request),
		})
	failOnError(err, "Failed to publish a message")
	if err == nil {
		updateJobStatus(jobID, "Queued")
	}
	log.Printf(" [x] Sent %s", request)
}
