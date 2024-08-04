package videoprocessor

import (
	"log"

	"github.com/Cirqach/GoStream/cmd/server/broadcast/videoProcessor/ffmpeg"
)

type VideoProcessor struct {

	FFmpegEngine *ffmpeg.FFmpeg

	
}

func NewVideoProcessor() *VideoProcessor{
	log.Println("Creating new video processor")
	return &VideoProcessor{
		FFmpegEngine: ffmpeg.NewFFmpeg(),
	}
}