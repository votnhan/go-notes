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

func (n *SQLiteLogNote) GetLastLogs(payload FetchLogsPayload) (logNotes []LogNote, nRows int64, err error) {
	err = n.db.Model(&LogNote{}).Count(&nRows).Error
	if err != nil {
		log.Printf("[GetLastLogs] - get number of rows error: %v", err)
		return
	}
	err = n.db.Limit(payload.Limit).Order("id desc").Find(&logNotes).Error
	if err != nil {
		log.Printf("[GetLastLogs] - get note logs error: %v", err)
		return
	}
	return
}

func (n *SQLiteLogNote) DeleteFirstLogs(payload RemoveLogsPayload) (rowsAffected int, err error) {
	var logIds []int
	err = n.db.Model(&LogNote{}).Select("Id").Limit(payload.NumberRows).Find(&logIds).Error
	if err != nil {
		log.Printf("[DeleteFirstLogs] - fetch ids error: %v", err)
		return
	}
	result := n.db.Where("Id in ?", logIds).Delete(&LogNote{})
	if result.Error != nil {
		log.Printf("[DeleteFirstLogs] - delete by ids error: %v", err)
		return
	}
	rowsAffected = int(result.RowsAffected)
	return
}
