package handler

import (
	"fmt"
	"time"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/model"
)

type VideoEncodingCompletedMessage struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	PublishedAt     string `json:"published_at"`
	Height          int    `json:"height"`
	Width           int    `json:"width"`
	DurationSeconds int    `json:"duration"`
	UserId          int    `json:"user_id"`
	OriginalId      string `json:"original_id"`
	Thumbnail       string `json:"thumbnail"`
	Path            string `json:"path"`
}

func ProcessVideoEncodingCompletedMessage(data *VideoEncodingCompletedMessage) error {
	publishedAt, err := time.Parse(time.RFC3339, data.PublishedAt)

	if err != nil {
		logger.Error("Invalid PublishedAt time format %v", err)

		return err
	}

	v := new(model.Video)

	if err = database.Conn.First(&v, "original_id = ?", data.OriginalId).Error; err == nil {
		return nil //row exists
	}

	v.Title = data.Title
	v.Description = data.Description
	v.OriginalId = data.OriginalId
	v.PublishedAt = publishedAt
	v.Duration = data.DurationSeconds
	v.Resolution = fmt.Sprintf("%dx%d", data.Width, data.Height)
	v.UserId = uint64(data.UserId)
	v.Thumbnail = data.Thumbnail
	v.Path = data.Path

	if err = database.Conn.Save(&v).Error; err != nil {
		logger.Error("Create video failed %v", err)

		return err
	}

	return nil
}
