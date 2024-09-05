package queuecontroller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Cirqach/GoStream/cmd/broadcast"
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/internal/database"
)

// QueueController struct    struct allow to control the video queue
type QueueController struct {
	dbController *database.DatabaseController
	c            chan database.Video // channel containing next video for broadcast
}

// NewQueueController function    create new queue controller
func NewQueueController(dbController *database.DatabaseController) *QueueController {
	return &QueueController{
		dbController: dbController,
		c:            make(chan database.Video),
	}
}

// TODO: fix 'http://localhost:8080/processed/./video/unprocessed/1.mp4/index.m3u8'); error
// StartControlling method    start controling schedule and update "broadcast"
func (q *QueueController) StartControlling(b *broadcast.Engine) {
	if err := q.dbController.ClearQueue(); err != nil {
		logger.LogError(logger.GetFuncName(0), "error clearing queue"+err.Error())
	}
	go q.controlSoonerVideo()
	go q.controlTime(b)

}

func (q *QueueController) controlTime(b *broadcast.Engine) {
	video := <-q.c

	if video.Time != time.Now().Format("2007-08-05 15:04:05") {
		b.UpdateVideo(video)
	}
	time.Sleep(1 * time.Second)
}

func (q *QueueController) controlSoonerVideo() {
	for {
		time.Sleep(5 * time.Second)
		video, err := q.dbController.GetSoonerVideo()
		if err != nil {
			logger.LogError(logger.GetFuncName(0), "error getting sooner video: "+err.Error())

		} else {
			err = q.dbController.DeleteSoonerVideo(video)
			if err != nil {
				logger.LogError(logger.GetFuncName(0), "error deleting sooner video: "+err.Error())
			}
		}
		q.c <- video
	}
}

func (q *QueueController) BookATime(wantedTime, date, filename, videoDuration string) error {
	// Parse the wanted time and date into a Go time object
	hour, err := strconv.Atoi(videoDuration[:2])
	if err != nil {
		return fmt.Errorf("invalid hour time format: %v", err)
	}
	minute, err := strconv.Atoi(videoDuration[3:5])
	if err != nil {
		return fmt.Errorf("invalid minute time format: %v", err)
	}
	second, err := strconv.Atoi(videoDuration[6:])
	if err != nil {
		return fmt.Errorf("invalid second time format: %v", err)
	}
	wantedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		logger.LogMessage(logger.GetFuncName(0), "Error parsing date:"+err.Error())
		return err
	}
	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Wanted date: %s", wantedDate))

	combinedTime := wantedDate.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute) + time.Second*time.Duration(second))

	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Combined time: %s", combinedTime))

	// Check if the wanted time overlaps with any existing broadcast times
	if q.dbController.CheckTimeOverlap(combinedTime, videoDuration) {
		return fmt.Errorf("Time is not free or overlap with other broadcast")
	}

	splitedTime := strings.Split(combinedTime.String(), " ")
	// If the time is free, add the video to the queue
	err = q.dbController.AddVideoToQueue(filename, splitedTime[0]+" "+splitedTime[1])
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return err
	}
	logger.LogMessage(logger.GetFuncName(0), "Video added to queue")
	return nil
}
