package queuecontroller

import (
	"time"

	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/internal/database"
)

// QueueController struct    struct allow to control the video queue
type QueueController struct {
	dbController *database.DatabaseController
}

// NewQueueController function    create new queue controller
func NewQueueController(dbController *database.DatabaseController) *QueueController {
	return &QueueController{
		dbController: dbController,
	}
}

// TODO: make it working, it not delete old queue and not change video on page
// StartControlling method    start controling schedule and update "broadcast"
func (q *QueueController) StartControlling(c chan database.Video, b *broadcast.Engine) {
	if err := q.dbController.ClearQueue(); err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
	}
	go q.controlSchedule(c)
	go b.UpdateVideo(c)
}

// controlSchedule method    controling video schedule
func (q *QueueController) controlSchedule(c chan database.Video) {
	for {
		// 1. Get the sooner video and its scheduled broadcast time
		video, err := q.dbController.GetSoonerVideo()
		if err != nil {
			logger.LogError(logger.GetFuncName(0), err.Error())
			continue // Move on to the next iteration in case of error
		}

		// 2. Calculate the duration until broadcast time
		now := time.Now()
		t, err := time.Parse(time.RFC3339, video.Time)
		if err != nil {
			logger.LogError(logger.GetFuncName(0), err.Error())
			continue
		}

		// Handle cases for broadcast time already passed or in the future
		if now.After(t) {
			// Broadcast time has already passed, send the video immediately
			c <- video
			continue
		}
		duration := t.Sub(now)

		// 3. Use time.After to create a timer for the remaining duration
		timer := time.After(duration)

		// 4. Select on the channel `c` and the timer for control flow
		select {
		case <-c:
			// Another video might have been pushed to the queue, handle it
			continue
		case <-timer:
			// Broadcast time reached, send the video for switching
			c <- video
		}
	}
}
