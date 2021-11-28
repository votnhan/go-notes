package main

import (
	"fmt"
	"golang/notes/db"
	"golang/notes/models"
	"golang/notes/tracking"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// Build rabbitmq publisher
	rabbitInfo := map[string]string{
		"Username":  Getenv("RABBIT_MQ_USER", "guest"),
		"Password":  Getenv("RABBIT_MQ_PASSWORD", "guest"),
		"Host":      Getenv("RABBIT_MQ_HOST", "rabbitmq"),
		"QueueName": Getenv("QUEUE_NAME", "log_note"),
	}
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v/", rabbitInfo["Username"], rabbitInfo["Password"], rabbitInfo["Host"]))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		rabbitInfo["QueueName"], // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// Build DB connection (writing log)
	db, _ := db.ConnectSQLite(Getenv("DB_FILE_LOG", "db/log_track.db"))
	logNoteModel := models.NewSQLiteLogNote(db)
	consumer := tracking.NewConsumer(queue, channel, logNoteModel)
	consumer.ConsumeMessage()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
