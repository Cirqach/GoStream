package ffmpeg

import (
	"fmt"
	"log"
	"os/exec"
)

type FFmpeg struct{

}
func NewFFmpeg() *FFmpeg{
	return &FFmpeg{}
}

func (f *FFmpeg) Process(inputFilePath, outputDirName string) error{
	log.Println("Processing video")
	mkdir(outputDirName)
	cmd := exec.Command("ffmpeg",
"-i", inputFilePath,
"-c:v", "libx264",
"-c:a", "aac",
"-hls_time", "10",
"-hls_list_size", "0",
// "-hls_segment_filename",  fmt.Sprintf("./video/processed/"+ "%s/segment_%03d.ts", outputDirName, ""),
"-hls_playlist_type", "vod",
fmt.Sprintf("%s/index.m3u8", "./video/processed/" + outputDirName))
log.Println("Running FFmpeg: " + cmd.String())
output, err := cmd.CombinedOutput()
if err != nil{
	log.Fatal(err)
	return fmt.Errorf("error running FFmpeg: %w\n%s", err, string(output))
}
fmt.Println("FFmpeg output: " + string(output))
return nil
}

func mkdir(outputDirName string){
	log.Println("Creating directory in path ./" + outputDirName)
cmd := exec.Command("mkdir",
"-p",
"./video/processed/" + outputDirName)
//	cmd := exec.Command("mkdir", "/video/processed/" + outputDirName)
	err := cmd.Run()
	if err != nil{
		log.Fatal(err)
	}
}
