package controllers

import (
	"golang/notes/models"
	"golang/notes/tracking"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"gorm.io/gorm"
)

type NoteCtrl struct {
	noteModel *models.SQLiteNote
	publisher *tracking.Publisher
	validate  *validator.Validate
}

func CreateNoteCtrl(noteModel *models.SQLiteNote, publisher *tracking.Publisher, validate *validator.Validate) *NoteCtrl {
	return &NoteCtrl{
		noteModel: noteModel,
		publisher: publisher,
		validate:  validate,
	}
}

func getInvalidParameterResponse(err error) map[string]interface{} {
	errMap := map[string]interface{}{
		"msg":    "Invalid parameter",
		"detail": strings.Split(string(err.Error()), "\n"),
	}
	return errMap
}

func (nc *NoteCtrl) CreateNote(c *gin.Context) {
	payload := models.CreateNotePayload{}
	if err := c.ShouldBind(&payload); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	noteUser := &models.Note{
		Content:   payload.Content,
		Title:     payload.Title,
		CreatedAt: time.Now().Round(0).String(),
		UpdatedAt: time.Now().Round(0).String(),
	}
	note, err := nc.noteModel.CreateNote(*noteUser)
	if err != nil {
		ResponseJSON(c, NewRestResponse(UNKNOWN_ERROR, nil, err.Error()), nil)
		return
	}
	go PushMessageTracking(nc.publisher, CREATE, nc.validate)
	ResponseJSON(c, NewRestResponse(SUCCESS, note, nil), nil)
}

func (nc *NoteCtrl) ReadNote(c *gin.Context) {
	uri_payload := models.IndentifyURIPayload{}
	if err := c.ShouldBindUri(&uri_payload); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	note, err := nc.noteModel.GetNote(uri_payload.NoteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, err.Error()), nil)
			return
		}
		ResponseJSON(c, NewRestResponse(UNKNOWN_ERROR, nil, err.Error()), nil)
		return
	}
	go PushMessageTracking(nc.publisher, READ, nc.validate)
	ResponseJSON(c, NewRestResponse(SUCCESS, note, nil), nil)
}

func (nc *NoteCtrl) ReadNoteAll(c *gin.Context) {
	payload := models.ReadAllNotePayload{
		Offset: 0,
		Limit:  5,
	}
	if err := c.ShouldBindWith(&payload, binding.Query); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	notes, err := nc.noteModel.GetNoteMany(payload)
	if err != nil {
		ResponseJSON(c, NewRestResponse(UNKNOWN_ERROR, nil, err.Error()), nil)
		return
	}
	go PushMessageTracking(nc.publisher, READ, nc.validate)
	ResponseJSON(c, NewRestResponse(SUCCESS, notes, nil), nil)
}

func (nc *NoteCtrl) UpdateNote(c *gin.Context) {
	uri_payload := models.IndentifyURIPayload{}
	if err := c.ShouldBindUri(&uri_payload); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	js_payload := models.UpdateNotePayload{}
	if err := c.ShouldBind(&js_payload); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	note, err := nc.noteModel.UpdateNote(js_payload, uri_payload.NoteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, err.Error()), nil)
			return
		}
		ResponseJSON(c, NewRestResponse(UNKNOWN_ERROR, nil, err.Error()), nil)
		return
	}
	go PushMessageTracking(nc.publisher, UPDATE, nc.validate)
	ResponseJSON(c, NewRestResponse(SUCCESS, note, nil), nil)
}

func (nc *NoteCtrl) DeleteNote(c *gin.Context) {
	uri_payload := models.IndentifyURIPayload{}
	if err := c.ShouldBindUri(&uri_payload); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	err := nc.noteModel.DeleteNote(uri_payload.NoteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, err.Error()), nil)
			return
		}
		ResponseJSON(c, NewRestResponse(UNKNOWN_ERROR, nil, err.Error()), nil)
		return
	}
	go PushMessageTracking(nc.publisher, DELETE, nc.validate)
	ResponseJSON(c, NewRestResponse(SUCCESS, "success", nil), nil)
}
