package handler

import (
	"fmt"
	"time"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/model"
)

type VideoEncodingCompleted struct {
	Title       string                             `json:"title"`
	Description string                             `json:"description"`
	PublishedAt string                             `json:"published_at"`
	Height      int                                `json:"height"`
	Width       int                                `json:"width"`
	Duration    int                                `json:"duration"`
	Resolutions []VideoEncodingCompletedResolution `json:"resolutions"`
	UserId      int                                `json:"user_id"`
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

	video := new(model.Video)

	video.Title = data.Title
	video.Description = data.Description
	video.Slug = fmt.Sprintf("%s-%s", helper.Slug(data.Title), helper.UniqueString(16))
	video.PublishedAt = publishedAt
	video.Duration = data.Duration
	video.Resolution = fmt.Sprintf("%dx%d", data.Width, data.Height)
	video.UserId = data.UserId

	if err = tx.Save(&video).Error; err != nil {
		log.Error("Create video failed %v", err)
		tx.Rollback()

		return err
	}

	for _, r := range data.Resolutions {
		for i, c := range r.Chunks {
			vc := new(model.VideoChunk)

			vc.VideoId = int(video.Id)
			vc.Order = i + 1
			vc.Resolution = fmt.Sprintf("%dx%d", r.Width, r.Height)
			vc.Encoding = r.Codec
			vc.Url = c

			if err = tx.Save(&vc).Error; err != nil {
				log.Error("Create chunk failed %v", err)
				tx.Rollback()

				return err
			}
		}
	}

	err = tx.Commit().Error

	if err != nil {
		return err
	}

	return nil
}
