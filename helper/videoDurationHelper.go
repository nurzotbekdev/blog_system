package helper

import (
	"encoding/json"
	"os/exec"
	"strconv"
)

type ffprobeResult struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func GetVideoDuration(path string) (int64, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		path,
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var result ffprobeResult
	if err := json.Unmarshal(out, &result); err != nil {
		return 0, err
	}

	sec, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, err
	}

	return int64(sec), nil
}
