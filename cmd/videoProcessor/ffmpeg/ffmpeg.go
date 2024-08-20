package ffmpeg

import (
	"fmt"
	"log"
	"os/exec"
)

// Parse function    parse video with ffmpeg
func Parse(inputFilePath, outputDirName string) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputFilePath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-hls_time", "10",
		"-hls_list_size", "0",
		// "-hls_segment_filename",  fmt.Sprintf("./video/processed/"+ "%s/segment_%03d.ts", outputDirName, ""),
		"-hls_playlist_type", "vod",
		fmt.Sprintf("%s/index.m3u8", "./video/processed/"+outputDirName))
	log.Println("Running FFmpeg: " + cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("error running FFmpeg: %w\n%s", err, string(output))
	}
	fmt.Println("FFmpeg output: " + string(output))
	return nil
}
