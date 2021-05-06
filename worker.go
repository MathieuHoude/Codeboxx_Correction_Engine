package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/getsentry/sentry-go"
	"github.com/streadway/amqp"
)

func worker(workerID int, queueName string) {
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
		defer sentry.Recover()
		for d := range msgs {
			var correctionRequest CorrectionRequest
			var gradingRequest GradingRequest
			var gradingResponse GradingResponse

			if queueName == "correction" {
				err = json.Unmarshal(d.Body, &correctionRequest)
				if err != nil {
					log.Println(err)
				}
				log.Printf(" [x] New "+queueName+" request for worker #"+fmt.Sprint(workerID)+": %s", correctionRequest.RepositoryURL)
				startCorrecting(correctionRequest)
			} else if queueName == "grading" {
				err = json.Unmarshal(d.Body, &gradingRequest)
				if err != nil {
					log.Println(err)
				}
				log.Printf(" [x] New "+queueName+" request for worker #"+fmt.Sprint(workerID)+": %s", gradingRequest.RepositoryURL)
				gradingResponse = startGrading(gradingRequest)
				updateJobStatus(gradingRequest.JobID, "Completed")
				sendBackResults(gradingRequest.JobID, gradingResponse)
			}

		}
	}()

	log.Printf(" [*] Waiting for " + queueName + " request.")
	<-forever
}
