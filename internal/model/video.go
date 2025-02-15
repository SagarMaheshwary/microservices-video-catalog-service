package model

import "time"

type Video struct {
	Id          uint64     `gorm:"primaryKey;autoIncrement;column:id"`
	Title       string     `gorm:"type:varchar(250);not null;column:title"`
	Description string     `gorm:"type:text;not null;column:description"`
	OriginalId  string     `gorm:"type:varchar(50);not null;column:original_id"`
	UserId      uint64     `gorm:"not null;column:user_id"`
	Resolution  string     `gorm:"type:varchar(25);not null;column:resolution"`
	Duration    int        `gorm:"not null;column:duration"`
	Path        string     `gorm:"type:varchar(250);not null;column:path"`
	Thumbnail   string     `gorm:"type:varchar(250);not null;column:thumbnail"`
	PublishedAt time.Time  `gorm:"not null;column:published_at"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (Video) TableName() string {
	return "videos"
}
