package queuecontroller

import (
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/internal/database"
)

type QueueController struct {
	dbController *database.DatabaseController
}

func NewQueueController(dbController *database.DatabaseController) *QueueController {
	return &QueueController{
		dbController: dbController,
	}
}

func (q *QueueController) StartControlling(c chan database.Video) {
	if err := q.dbController.ClearQueue(); err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
	}
	go q.controlSchedule(c)
}

func (q *QueueController) controlSchedule(c chan database.Video) {
	for {
		video, err := q.dbController.GetSoonerVideo()
		if err != nil {
			logger.LogError(logger.GetFuncName(0), err.Error())
		}
		c <- video
	}
}
