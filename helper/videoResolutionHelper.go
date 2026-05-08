package helper

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type ffprobeStream struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func GetResolution(path string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "json",
		path,
	)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var result ffprobeStream
	if err := json.Unmarshal(out, &result); err != nil {
		return "", err
	}

	if len(result.Streams) == 0 {
		return "", fmt.Errorf("video stream topilmadi")
	}

	w := result.Streams[0].Width
	h := result.Streams[0].Height

	return strconv.Itoa(w) + "x" + strconv.Itoa(h), nil
}
