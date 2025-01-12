package handler

import (
	"fmt"
	"time"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/model"
)

type VideoEncodingCompletedMessage struct {
	Title              string              `json:"title"`
	Description        string              `json:"description"`
	PublishedAt        string              `json:"published_at"`
	Height             int                 `json:"height"`
	Width              int                 `json:"width"`
	Duration           int                 `json:"duration"`
	EncodedResolutions []EncodedResolution `json:"resolutions"`
	UserId             int                 `json:"user_id"`
	OriginalId         string              `json:"original_id"`
	Thumbnail          string              `json:"thumbnail"`
	Path               string              `json:"path"`
}

type EncodedResolution struct {
	Height int      `json:"height"`
	Width  int      `json:"width"`
	Chunks []string `json:"chunks"`
}

func ProcessVideoEncodingCompletedMessage(data *VideoEncodingCompletedMessage) error {
	publishedAt, err := time.Parse(time.RFC3339, data.PublishedAt)

	if err != nil {
		logger.Error("Invalid PublishedAt time format %v", err)

		return err
	}

	tx := database.Conn.Begin()

	v := new(model.Video)

	res := tx.First(&v, "original_id = ?", data.OriginalId)

	if res.Error != nil {
		logger.Info("Processing video row.")

		v.Title = data.Title
		v.Description = data.Description
		v.OriginalId = data.OriginalId
		v.PublishedAt = publishedAt
		v.Duration = data.Duration
		v.Resolution = fmt.Sprintf("%dx%d", data.Width, data.Height)
		v.UserId = data.UserId
		v.Thumbnail = data.Thumbnail
		v.Path = data.Path

		if err = tx.Save(&v).Error; err != nil {
			logger.Error("Create video failed %v", err)
			tx.Rollback()

			return err
		}
	}

	logger.Info("Processing video_chunk rows.")

	for _, r := range data.EncodedResolutions {
		for _, c := range r.Chunks {
			vc := new(model.VideoChunk)
			resolution := fmt.Sprintf("%dx%d", r.Width, r.Height)

			res := tx.Where(map[string]interface{}{"video_id": v.Id, "resolution": resolution, "filename": c}).First(&vc)

			if res.Error != nil {
				vc.VideoId = int(v.Id)
				vc.Resolution = resolution
				vc.Filename = c

				if err := tx.Save(&vc).Error; err != nil {
					logger.Error("Create video_chunk failed %v", err)
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

	logger.Info("Transaction committed.")

	return nil
}
