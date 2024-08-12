package videoprocessor

import (
	"log"
	"os/exec"

	"github.com/Cirqach/GoStream/cmd/videoProcessor/ffmpeg"
)

type VideoProcessor struct {
	FFmpegEngine *ffmpeg.FFmpeg
}

func NewVideoProcessor() *VideoProcessor {
	log.Println("Creating new video processor")
	return &VideoProcessor{
		FFmpegEngine: ffmpeg.NewFFmpeg(),
	}
}

func (vp *VideoProcessor) Process(inputFilePath, outputDirName string) {
	mkdir(outputDirName)
	vp.FFmpegEngine.Parse(inputFilePath, outputDirName)

}

func mkdir(outputDirName string) {
	log.Println("Creating directory in path /video/processed/" + outputDirName)
	cmd := exec.Command("mkdir",
		"-p",
		"./video/processed/"+outputDirName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func delete(dirName string) {
	log.Println("Deleting directory in path /video/processed/" + dirName)
	cmd := exec.Command("rm",
		"-rf",
		"./video/processed/"+dirName)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
