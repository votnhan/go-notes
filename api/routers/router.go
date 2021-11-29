package routers

import (
	"fmt"
	"golang/notes/api/controllers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type RabbitMQSecret struct {
	Username  string
	Password  string
	Host      string
	QueueName string
}

func HealthCheck(c *gin.Context) {
	c.JSON(controllers.HttpStatusMap[controllers.SUCCESS], map[string]string{})
}

func CreateRouter(noteCtrl *controllers.NoteCtrl, logNoteCtrl *controllers.LogNoteCtrl) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	// for healthcheck
	v1.GET("/healthcheck", HealthCheck)
	// for CRUD notes
	v1.GET("/note/:noteid", noteCtrl.ReadNote)
	v1.PUT("/note/:noteid", noteCtrl.UpdateNote)
	v1.DELETE("/note", noteCtrl.DeleteNote)
	v1.POST("/note", noteCtrl.CreateNote)
	v1.GET("/note", noteCtrl.ReadNoteAll)
	// for tracking logs
	v1.GET("/logs", logNoteCtrl.FetchLogNotes)
	v1.DELETE("/logs", logNoteCtrl.DeleteLogNotes)
	return router
}

func BuildQueue(queueSecret *RabbitMQSecret) (conn *amqp.Connection, queue amqp.Queue, channel *amqp.Channel, err error) {
	conn, err = amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v/", queueSecret.Username, queueSecret.Password, queueSecret.Host))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}
	channel, err = conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
	}
	queue, err = channel.QueueDeclare(
		queueSecret.QueueName, // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare a queue")
	}
	return
}
