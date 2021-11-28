package models

import (
	"log"

	"gorm.io/gorm"
)

type SQLiteLogNote struct {
	db *gorm.DB
}

func NewSQLiteLogNote(db *gorm.DB) *SQLiteLogNote {
	return &SQLiteLogNote{
		db: db,
	}
}

type LogNote struct {
	Id        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskType  string `gorm:"C|R|U|D;not null" json:"task_type"`
	CreatedAt string `gorm:"not null" json:"created_at"`
}

type FilteredIdLogNote struct {
	Id int `gorm:"primaryKey;autoIncrement" json:"id"`
}

type FetchLogsPayload struct {
	Limit int `binding:"gte=0" json:"limit" form:"limit"`
}

type RemoveLogsPayload struct {
	NumberRows int `binding:"required,gte=1" json:"rows" form:"rows"`
}

func (LogNote) TableName() string {
	return "log_note"
}

func (n *SQLiteLogNote) CreateLogNote(logNoteInput LogNote) (logNote LogNote, err error) {
	err = n.db.Create(&logNoteInput).Error
	if err != nil {
		log.Printf("[CreateLogNote] - error: %v", err)
		return
	}
	err = n.db.First(&logNote, "id = ?", logNoteInput.Id).Error
	if err != nil {
		log.Printf("[CreateLogNote] - re-query error: %v", err)
		return
	}
	return
}

func (n *SQLiteLogNote) GetLastLogs(payload FetchLogsPayload) (logNotes []LogNote, err error) {
	err = n.db.Limit(payload.Limit).Order("id desc").Find(&logNotes).Error
	if err != nil {
		log.Printf("[GetLastLogs] - error: %v", err)
		return
	}
	return
}

func (n *SQLiteLogNote) DeleteFirstLogs(payload RemoveLogsPayload) (err error) {
	var logIds []FilteredIdLogNote
	err = n.db.Model(&LogNote{}).Limit(payload.NumberRows).Find(&logIds).Error
	if err != nil {
		log.Printf("[DeleteFirstLogs] - fetch id error: %v", err)
		return
	}
	err = n.db.Model(&LogNote{}).Delete(&logIds).Error
	if err != nil {
		log.Printf("[DeleteFirstLogs] - delete by ids error: %v", err)
		return
	}
	return
}
