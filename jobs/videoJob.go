package jobs

type VideoJob struct {
	VideoID    uint   `json:"video_id"`
	FilePath   string `json:"file_path"`
	Resolution string `json:"resolution"`
}
