package queuecontroller

import (
	"fmt"
	"strings"
	"time"

	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/internal/database"
)

// QueueController struct    struct allow to control the video queue
type QueueController struct {
	dbController *database.DatabaseController
	c            chan database.Video
}

// NewQueueController function    create new queue controller
func NewQueueController(dbController *database.DatabaseController) *QueueController {
	return &QueueController{
		dbController: dbController,
		c:            make(chan database.Video),
	}
}

// StartControlling method    start controling schedule and update "broadcast"
func (q *QueueController) StartControlling(b *broadcast.Engine) {
	if err := q.dbController.ClearQueue(); err != nil {
		logger.LogError(logger.GetFuncName(0), "error clearing queue"+err.Error())
	}
	go q.controlSchedule()
	go b.UpdateVideo(q.c)
}

// TODO: fix infinite loop
// controlSchedule method    controling video schedule
func (q *QueueController) controlSchedule() {
	for {
		// 1. Get the sooner video and its scheduled broadcast time
		video, err := q.dbController.GetSoonerVideo()
		if err != nil {

			logger.LogError(logger.GetFuncName(0), err.Error())
			time.Sleep(time.Second * 10)
			continue // Move on to the next iteration in case of error
		}

		if video == (database.Video{}) {
			logger.LogMessage(logger.GetFuncName(0), "No videos in queue")
			time.Sleep(time.Second * 10)
			continue
		}

		// 2. Calculate the duration until broadcast time
		now := time.Now()
		t, err := time.Parse(time.RFC3339, video.Time)
		if err != nil {
			logger.LogError(logger.GetFuncName(0), err.Error())
			time.Sleep(time.Second * 10)
			continue
		}

		// Handle cases for broadcast time already passed or in the future
		if now.After(t) {
			// Broadcast time has already passed, send the video immediately
			q.c <- video
			continue
		}
		duration := t.Sub(now)

		// 3. Use time.After to create a timer for the remaining duration
		timer := time.After(duration)

		// 4. Select on the channel `c` and the timer for control flow
		select {
		case <-q.c:
			// Another video might have been pushed to the queue, handle it
			continue
		case <-timer:
			// Broadcast time reached, send the video for switching
			q.c <- video
		}
	}
}

// TODO: add good logic for booking time
func (q *QueueController) BookATime(wantedTime, date, filename, videoDuration string) error {
	// Parse the wanted time and date into a Go time object
	wantedTimeObj, err := time.Parse("HH:mm:ss", wantedTime)
	if err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}
	month, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}

	wantedDateTime := time.Date(
		time.Now().Year(),
		month.Month(),
		day.Day(),
		wantedTimeObj.Hour(),
		wantedTimeObj.Minute(),
		wantedTimeObj.Second(),
		0,
		time.UTC)

	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Wanted time: %s, date: %s", wantedDateTime, date))
	// Check if the wanted time overlaps with any existing broadcast times
	if !q.dbController.CheckTimeOverlap(wantedDateTime, videoDuration) {
		return fmt.Errorf("Time is not free or overlap with other broadcast")
	}

	// If the time is free, add the video to the queue
	err = q.dbController.AddVideoToQueue(filename, strings.Split(wantedDateTime.String(), ".")[0])
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return err
	}
	logger.LogMessage(logger.GetFuncName(0), "Video added to queue")
	return nil
}
