package model

import "time"

type VideoChunk struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	VideoId    int       `json:"video_id"`
	Resolution string    `json:"resolution"`
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
