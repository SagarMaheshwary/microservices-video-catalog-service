package handler

import (
	"fmt"
	"time"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/model"
)

type VideoEncodingCompleted struct {
	Title        string                             `json:"title"`
	Description  string                             `json:"description"`
	PublishedAt  string                             `json:"published_at"`
	Height       int                                `json:"height"`
	Width        int                                `json:"width"`
	Duration     int                                `json:"duration"`
	Resolutions  []VideoEncodingCompletedResolution `json:"resolutions"`
	UserId       int                                `json:"user_id"`
	OriginalId   string                             `json:"original_id"`
	ThumbnailUrl string                             `json:"thumbnail_url"`
}

type VideoEncodingCompletedResolution struct {
	Height int      `json:"height"`
	Width  int      `json:"width"`
	Codec  string   `json:"codec"`
	Chunks []string `json:"chunks"`
}

func ProcessVideoEncodingCompleted(data *VideoEncodingCompleted) error {
	publishedAt, err := time.Parse(time.RFC3339, data.PublishedAt)

	if err != nil {
		log.Error("Invalid PublishedAt time format %v", err)

		return err
	}

	tx := database.DB.Begin()

	v := new(model.Video)

	res := tx.First(&v, "original_id = ?", data.OriginalId)

	if res.Error != nil {
		v.Title = data.Title
		v.Description = data.Description
		v.OriginalId = data.OriginalId
		v.PublishedAt = publishedAt
		v.Duration = data.Duration
		v.Resolution = fmt.Sprintf("%dx%d", data.Width, data.Height)
		v.UserId = data.UserId
		v.ThumbnailUrl = data.ThumbnailUrl

		if err = tx.Save(&v).Error; err != nil {
			log.Error("Create video failed %v", err)
			tx.Rollback()

			return err
		}

	}

	for _, r := range data.Resolutions {
		for i, c := range r.Chunks {
			vc := new(model.VideoChunk)
			order := i + 1
			resolution := fmt.Sprintf("%dx%d", r.Width, r.Height)

			res := tx.Where(map[string]interface{}{"video_id": v.Id, "order": order, "resolution": resolution}).First(&vc)

			if res.Error != nil {
				vc.VideoId = int(v.Id)
				vc.Order = order
				vc.Resolution = resolution
				vc.Encoding = r.Codec
				vc.Url = c

				if err := tx.Save(&vc).Error; err != nil {
					log.Error("Create chunk failed %v", err)
					tx.Rollback()

					return err
				}
			}
		}
	}

	err = tx.Commit().Error

	if err != nil {
		return err
	}

	return nil
}
