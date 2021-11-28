package tracking

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"

	"golang/notes/models"
	"time"

	"github.com/streadway/amqp"
)

type Consumer struct {
	Queue   amqp.Queue
	Channel *amqp.Channel
	Model   *models.SQLiteLogNote
}

func NewConsumer(queue amqp.Queue, channel *amqp.Channel, model *models.SQLiteLogNote) *Consumer {
	return &Consumer{
		Queue:   queue,
		Channel: channel,
		Model:   model,
	}
}

func FromGOB64(str string) Message {
	m := Message{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Printf("failed base64 Decode: %v", err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		log.Printf("failed gob Decode: %v", err)
	}
	return m
}

func (c *Consumer) SaveMessage(taskType string) (err error) {
	logNoteInput := models.LogNote{
		TaskType:  taskType,
		CreatedAt: time.Now().Round(0).String(),
	}
	_, err = c.Model.CreateLogNote(logNoteInput)
	if err != nil {
		log.Fatalf("Fail in saving log: %v", err)
	}
	return
}

func (c *Consumer) ConsumeMessage() {
	messages, err := c.Channel.Consume(
		c.Queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer")
	}
	forever := make(chan bool)
	go func() {
		for msg := range messages {
			body := FromGOB64(string(msg.Body))
			err = c.SaveMessage(body["task_type"].(string))
			log.Printf("Writing log for task: %s", body["task_type"].(string))
			msg.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
