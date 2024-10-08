package ffmpeg

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/Cirqach/GoStream/cmd/logger"
)

// Parse function    parse video with given name from path video/unprocessed/ with ffmpeg
func Parse(fileName string) error {
	logger.LogMessage(logger.GetFuncName(0), "Parsing video with ffmpeg")
	cmd := exec.Command("ffmpeg",
		"-i", fileName,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-hls_time", "10",
		"-hls_list_size", "0",
		// "-hls_segment_filename",  fmt.Sprintf("./video/processed/"+ "%s/segment_%03d.ts", outputDirName, ""),
		"-hls_playlist_type", "vod",
		fmt.Sprintf("%s/index.m3u8", "./video/processed/"+fileName))
	log.Println("Running FFmpeg: " + cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("error running FFmpeg: %w\n%s", err, string(output))
	}
	fmt.Println("FFmpeg output: " + string(output))
	return nil
}

func GetVideoDuration(filePath string) (time.Time, error) {
	// using ffprobe for getting video duration
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", "./unprocessed")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogError(logger.GetFuncName(0), "Error getting video duration:"+err.Error())
		return time.Time{}, err
	}

	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Video duration: %s", string(output)))
	// deleting milliseconds from output
	duration, err := time.Parse("15:04:05", strings.Split(string(output), ".")[0])
	if err != nil {
		logger.LogError(logger.GetFuncName(0), "Error parsing video duration:"+err.Error())
		return time.Time{}, err
	}

	logger.LogMessage(logger.GetFuncName(0), fmt.Sprintf("Video duration: %s", duration.String()))
	// returning video duration
	return duration, nil
}
