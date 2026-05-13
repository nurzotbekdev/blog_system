package helper

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func ConvertVideo(inputPath string, quality string) (string, int64, string, error) {
	var height string

	switch quality {
	case "1080p":
		height = "1080"
	case "720p":
		height = "720"
	case "480p":
		height = "480"
	case "360p":
		height = "360"
	default:
		return "", 0, "", fmt.Errorf("unknown quality")
	}

	outputPath := fmt.Sprintf(
		"uploads/video/%s_%s.mp4",
		uuid.New().String(),
		quality,
	)

	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-vf", "scale=-2:"+height,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "aac",
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return "", 0, "", err
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		return "", 0, "", err
	}

	return outputPath, info.Size(), "mp4", nil
}
