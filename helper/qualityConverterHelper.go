package helper

import "fmt"

func QualityConverter(resolution string) []string {
	var w, h int

	fmt.Sscanf(resolution, "%dx%d", &w, &h)

	switch {
	case h >= 2160:
		return []string{"1080p", "720p", "480p", "360p"}
	case h >= 1440:
		return []string{"1080p", "720p", "480p", "360p"}
	case h >= 1080:
		return []string{"720p", "480p", "360p"}
	case h >= 720:
		return []string{"480p", "360p"}
	case h >= 480:
		return []string{"360p"}
	default:
		return []string{}
	}
}
