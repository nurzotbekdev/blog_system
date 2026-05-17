package workers

import (
	"blog_system/config"
	"blog_system/jobs"
	"blog_system/models"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

func StartVideoViewWorker() {

	for {

		result, err := config.RedisClient.BRPop(
			config.Ctx,
			0,
			"video_view_queue",
		).Result()

		if err != nil {
			continue
		}

		var job jobs.VideoViewJob
		if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
			continue
		}

		process(job)
	}
}

func process(job jobs.VideoViewJob) {
	config.DB.Model(&models.Video{}).
		Where("id = ?", job.VideoID).
		Update("views", gorm.Expr("views + 1"))

	config.DB.Model(&models.Channel{}).
		Where("id = (?)",
			config.DB.Table("videos").
				Select("channel_id").
				Where("id = ?", job.VideoID),
		).
		Update("total_views", gorm.Expr("total_views + 1"))

	var history models.History

	err := config.DB.
		Where("user_id = ? AND video_id = ?", job.UserID, job.VideoID).
		First(&history).Error

	if err == nil {
		config.DB.Model(&history).
			Update("updated_at", time.Now())
	} else {
		config.DB.Create(&models.History{
			UserID:  job.UserID,
			VideoID: job.VideoID,
		})
	}
}
