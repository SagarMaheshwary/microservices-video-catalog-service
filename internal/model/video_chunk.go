package model

import "time"

type VideoChunk struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	VideoId    int       `json:"video_id"`
	Order      int       `json:"order"`
	Resolution string    `json:"resolution"`
	Encoding   string    `json:"encoding"`
	Url        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
