package workers

import (
	"blog_system/config"
	"blog_system/helper"
	"blog_system/jobs"
	"blog_system/models"
	"context"
	"encoding/json"
	"log"
)

func StartVideoWorker() {
	for {
		result, err := config.RedisClient.BRPop(context.Background(), 0, "video_queue").Result()
		if err != nil {
			log.Println("Redis error:", err)
			continue
		}

		var job jobs.VideoJob
		if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
			log.Println("Invalid job:", err)
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Worker panic recovered:", r)
				}
			}()

			processVideo(job)
		}()
	}
}

func processVideo(job jobs.VideoJob) {
	qualities := helper.QualityConverter(job.Resolution)

	for _, q := range qualities {

		url, size, format, err := helper.ConvertVideo(job.FilePath, q)
		if err != nil {
			log.Println("convert error:", err)
			continue
		}

		var videoQuality models.VideoQuality

		err = config.DB.FirstOrCreate(&videoQuality,
			models.VideoQuality{
				VideoID: job.VideoID,
				Quality: q,
			},
		).Error

		if err != nil {
			log.Println("db insert failed:", err)
			continue
		}

		videoQuality.VideoURL = url
		videoQuality.Size = size
		videoQuality.Format = format

		if err := config.DB.Save(&videoQuality).Error; err != nil {
			log.Println("db update failed:", err)
			continue
		}
	}

	if err := config.DB.Model(&models.Video{}).
		Where("id = ?", job.VideoID).
		Update("status", "ready").Error; err != nil {
		log.Println("status update failed:", err)
	}
}
