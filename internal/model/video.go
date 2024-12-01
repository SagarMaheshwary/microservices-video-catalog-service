package model

import "time"

type Video struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	UserId      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OriginalId  string    `json:"original_id"`
	Resolution  string    `json:"resolution"`
	Duration    int       `json:"duration"`
	Path        string    `json:"path"`
	Thumbnail   string    `json:"thumbnail"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
