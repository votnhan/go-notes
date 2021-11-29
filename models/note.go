package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type SQLiteNote struct {
	db *gorm.DB
}

func NewSQLiteNote(db *gorm.DB) *SQLiteNote {
	return &SQLiteNote{
		db: db,
	}
}

type Note struct {
	Id        int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string `gorm:"size:512;not null" json:"title"`
	Content   string `gorm:"size:1024;not null" json:"content"`
	CreatedAt string `gorm:"not null" json:"created_at"`
	UpdatedAt string `gorm:"not null" json:"updated_at"`
}

type CreateNotePayload struct {
	Title   string `binding:"required,lte=512" json:"title"`
	Content string `binding:"required,lte=1024" json:"content"`
}

type UpdateNotePayload struct {
	Title   string `binding:"lte=512" json:"title"`
	Content string `binding:"lte=1024" json:"content"`
}

type IndentifyURIPayload struct {
	NoteId int `binding:"required" uri:"noteid"`
}

type DeleteNotePayload struct {
	NoteIds []int `binding:"required,min=1" json:"note_ids"`
}

type ReadAllNotePayload struct {
	Offset int `binding:"gte=0" json:"offset" form:"offset"`
	Limit  int `binding:"gte=0" json:"limit" form:"limit"`
}

func (Note) TableName() string {
	return "note"
}

func (n *SQLiteNote) CreateNote(noteUser Note) (note Note, err error) {
	err = n.db.Create(&noteUser).Error
	if err != nil {
		log.Printf("[CreateNote] - error: %v", err)
		return
	}
	err = n.db.First(&note, "id = ?", noteUser.Id).Error
	if err != nil {
		log.Printf("[CreateNote] - re-query error: %v", err)
		return
	}
	return
}

func (n *SQLiteNote) GetNote(noteId int) (note Note, err error) {
	err = n.db.First(&note, "id = ?", noteId).Error
	if err != nil {
		log.Printf("[GetNote] - error: %v", err)
		return
	}
	return
}

func (n *SQLiteNote) GetNoteMany(payload ReadAllNotePayload) (notes []Note, nRows int64, err error) {
	err = n.db.Model(&Note{}).Count(&nRows).Error
	if err != nil {
		log.Printf("[GetNoteMany] - get number of rows error: %v", err)
		return
	}
	err = n.db.Limit(payload.Limit).Offset(payload.Offset).Order("id desc").Find(&notes).Error
	if err != nil {
		log.Printf("[GetNoteMany] - get notes error: %v", err)
		return
	}
	return
}

func (n *SQLiteNote) UpdateNote(payload UpdateNotePayload, noteId int) (note Note, err error) {
	noteUser := Note{
		Id:        noteId,
		Title:     payload.Title,
		Content:   payload.Content,
		UpdatedAt: time.Now().Round(0).String(),
	}
	err = n.db.Updates(noteUser).Error
	if err != nil {
		log.Printf("[UpdateNote] - error: %v", err)
		return
	}
	err = n.db.First(&note, "id = ?", noteUser.Id).Error
	if err != nil {
		log.Printf("[UpdateNote] - re-query error: %v", err)
		return
	}
	return
}

func (n *SQLiteNote) DeleteNote(payload DeleteNotePayload) (rowsAffected int, err error) {
	deleteResult := n.db.Delete(&Note{}, payload.NoteIds)
	if deleteResult.Error != nil {
		log.Printf("[DeleteNote] - error: %v", err)
		return 0, deleteResult.Error
	}
	rowsAffected = int(deleteResult.RowsAffected)
	return
}
