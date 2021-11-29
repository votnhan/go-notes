package controllers

import (
	"golang/notes/models"

	"github.com/gin-gonic/gin"
)

type LogNoteCtrl struct {
	LogNoteModel *models.SQLiteLogNote
}

func CreateLogNoteCtrl(logNoteModel *models.SQLiteLogNote) *LogNoteCtrl {
	return &LogNoteCtrl{
		LogNoteModel: logNoteModel,
	}
}

func (lc *LogNoteCtrl) FetchLogNotes(c *gin.Context) {
	payload := models.FetchLogsPayload{
		Limit: 5,
	}
	if err := c.ShouldBind(&payload); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	logNotes, nRows, err := lc.LogNoteModel.GetLastLogs(payload)
	if err != nil {
		ResponseJSON(c, NewRestResponse(UNKNOWN_ERROR, nil, err.Error()), nil)
		return
	}
	response := map[string]interface{}{
		"all":   nRows,
		"fetch": len(logNotes),
		"data":  logNotes,
	}
	ResponseJSON(c, NewRestResponse(SUCCESS, response, nil), nil)
}

func (lc *LogNoteCtrl) DeleteLogNotes(c *gin.Context) {
	payload := models.RemoveLogsPayload{}
	if err := c.ShouldBind(&payload); err != nil {
		errMap := getInvalidParameterResponse(err)
		ResponseJSON(c, NewRestResponse(INVALID_PARAMETERS, nil, errMap), nil)
		return
	}
	rowsAffected, err := lc.LogNoteModel.DeleteFirstLogs(payload)
	if err != nil {
		ResponseJSON(c, NewRestResponse(UNKNOWN_ERROR, nil, err.Error()), nil)
		return
	}
	response := map[string]interface{}{
		"rows_affected": rowsAffected,
	}
	ResponseJSON(c, NewRestResponse(SUCCESS, response, nil), nil)
}
