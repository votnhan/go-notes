package main

import (
	"golang/notes/api/controllers"
	"golang/notes/api/routers"
	"golang/notes/db"
	"golang/notes/models"
	"golang/notes/tracking"
	"os"

	"gopkg.in/go-playground/validator.v8"
)

func main() {
	// Build rabbitmq publisher
	rabbitmqInfo := &routers.RabbitMQSecret{
		Username:  Getenv("RABBIT_MQ_USER", "guest"),
		Password:  Getenv("RABBIT_MQ_PASSWORD", "guest"),
		Host:      Getenv("RABBIT_MQ_HOST", "rabbitmq"),
		QueueName: Getenv("QUEUE_NAME", "log_note"),
	}
	connection, queue, channel, _ := routers.BuildQueue(rabbitmqInfo)
	publisher := tracking.NewPublisher(queue, channel)
	defer connection.Close()
	defer channel.Close()

	// Build DB connection
	dbNote, _ := db.ConnectSQLite(Getenv("DB_FILE_CRUD", "db/go_note.db"))
	noteModel := models.NewSQLiteNote(dbNote)
	dbLogNote, _ := db.ConnectSQLite(Getenv("DB_FILE_LOG", "db/log_track.db"))
	logNoteModel := models.NewSQLiteLogNote(dbLogNote)

	// Build validator
	validate := validator.New(&validator.Config{TagName: "validate"})

	// Build controllers
	noteCtrl := controllers.CreateNoteCtrl(noteModel, publisher, validate)
	logNoteCtrl := controllers.CreateLogNoteCtrl(logNoteModel)

	// Build routers
	router := routers.CreateRouter(noteCtrl, logNoteCtrl)

	// Run app
	router.Run(os.Getenv("HOST"))
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
