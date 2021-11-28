package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"golang/notes/tracking"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

type ResponseCode int
type TaskType int
type ResponseResult interface{}
type ResponseError interface{}

const (
	SUCCESS ResponseCode = 0

	// Common service errors (1-10)
	UNKNOWN_ERROR      ResponseCode = 1
	UNAUTHORIZED       ResponseCode = 2
	PERMISSION_DENIED  ResponseCode = 3
	INVALID_PARAMETERS ResponseCode = 4

	// Data errors
	RECORD_DOES_NOT_EXIST ResponseCode = 11
	RECORD_EXISTS         ResponseCode = 12
	RECORD_HAS_NO_CHANGE  ResponseCode = 13
)

const (
	CREATE TaskType = 0
	READ   TaskType = 1
	UPDATE TaskType = 2
	DELETE TaskType = 3
)

var (
	HttpStatusMap = map[ResponseCode]int{
		SUCCESS:               http.StatusOK,
		UNKNOWN_ERROR:         http.StatusInternalServerError,
		UNAUTHORIZED:          http.StatusUnauthorized,
		PERMISSION_DENIED:     http.StatusInternalServerError,
		INVALID_PARAMETERS:    http.StatusBadRequest,
		RECORD_DOES_NOT_EXIST: http.StatusNotFound,
		RECORD_EXISTS:         http.StatusConflict,
		RECORD_HAS_NO_CHANGE:  http.StatusConflict,
	}
)

var (
	TaskTypeMap = map[TaskType]string{
		CREATE: "C",
		READ:   "R",
		UPDATE: "U",
		DELETE: "D",
	}
)

type RestResponse struct {
	Code   ResponseCode   `json:"code"`
	Result ResponseResult `json:"result,omitempty"`
	Error  ResponseError  `json:"error,omitempty"`
}

func NewRestResponse(code ResponseCode, result ResponseResult, err ResponseError) *RestResponse {
	r := RestResponse{
		Code:   code,
		Result: result,
		Error:  err,
	}
	return &r
}

func ResponseJSON(c *gin.Context, r *RestResponse, status *int) error {
	_, err := json.Marshal(r)
	if err != nil {
		return err
	}
	c.Header("Content-Type", "application/json")
	if status == nil {
		s, ok := HttpStatusMap[r.Code]
		if !ok {
			c.JSON(http.StatusInternalServerError, structs.Map(r))
		} else {
			c.JSON(s, structs.Map(r))
		}
		return nil
	}
	c.JSON(*status, structs.Map(r))
	return nil
}

func PushMessageTracking(publisher *tracking.Publisher, taskType TaskType, validate *validator.Validate) (err error) {
	payload := &tracking.CreateMessagePayload{
		TaskType: TaskTypeMap[taskType],
	}
	err = validate.Struct(payload)
	if err != nil {
		log.Printf("Invalid message: %v", err)
	}
	err = publisher.PushMessage(payload)
	return
}
