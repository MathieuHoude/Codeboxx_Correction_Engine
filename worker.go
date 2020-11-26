package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func worker(workerID int) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var gradingRequest GradingRequest
			err = json.Unmarshal(d.Body, &gradingRequest)
			if err != nil {
				log.Println(err)
			}
			log.Printf(" [x] New grading request for worker #"+fmt.Sprint(workerID)+": %s", gradingRequest.ProjectName)

			gradingResponse := startGrading(gradingRequest)
			jsonData, err := json.Marshal(gradingResponse)
			failOnError(err, "Failed to parse struct to JSON")

			fmt.Println(string(jsonData))
		}
	}()

	log.Printf(" [*] Waiting for grading request.")
	<-forever
}
