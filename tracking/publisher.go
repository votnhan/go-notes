package tracking

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"

	"github.com/streadway/amqp"
)

type Publisher struct {
	Queue   amqp.Queue
	Channel *amqp.Channel
}

type Message map[string]interface{}

type CreateMessagePayload struct {
	TaskType string `validate:"required" json:"task_type"`
}

func NewPublisher(queue amqp.Queue, channel *amqp.Channel) *Publisher {
	return &Publisher{
		Queue:   queue,
		Channel: channel,
	}
}

func ToGOB64(message Message) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(message)
	if err != nil {
		log.Printf("Failed gob Encode: %v", err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (p *Publisher) PushMessage(payload *CreateMessagePayload) (err error) {
	message := Message{
		"task_type": payload.TaskType,
	}
	encodedMsg := ToGOB64(message)
	err = p.Channel.Publish(
		"",
		p.Queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(encodedMsg),
		})
	if err != nil {
		log.Printf("[PublishMessage] - error: %v", err)
		return
	}
	return
}
