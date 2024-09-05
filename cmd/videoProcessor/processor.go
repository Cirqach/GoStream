package videoprocessor

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"

	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/Cirqach/GoStream/cmd/videoProcessor/ffmpeg"
)

// Process function    process video from /video/unprocessed/ with given name
func Process(fileName string) error {
	mkdir(fileName)
	err := ffmpeg.Parse(fileName)
	return err

}

// mkdir function    create directory in path /video/processed/ with given name
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

// delete function    delete directory in path /video/processed/ with given name
// func delete(dirName string) {
// 	log.Println("Deleting directory in path /video/processed/" + dirName)
// 	cmd := exec.Command("rm",
// 		"-rf",
// 		"./video/processed/"+dirName)
// 	err := cmd.Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// SaveVideo function    save video in path /video/unprocessed/ with given name
func SaveVideo(file multipart.File, filename string) error {
	logger.LogMessage(logger.GetFuncName(0), "Saving video in path /video/unprocessed/"+filename)
	dst, err := os.Create("./video/unprocessed/" + filename)
	if err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		logger.LogError(logger.GetFuncName(0), err.Error())
		return err
	}
	logger.LogMessage(logger.GetFuncName(0), "Video saved in path /video/unprocessed/"+filename)
	return nil
}
