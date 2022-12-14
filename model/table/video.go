package table

import "time"

type Video struct {
	ID               int64     `gorm:"column:id" json:"id"`
	VideoName        string    `gorm:"column:video_name" json:"video_name"`
	VideoDescription string    `gorm:"column:video_description" json:"video_description"`
	MagnetLink       string    `gorm:"column:magnet_link" json:"magnet_link"`
	VideoSize        int64     `gorm:"column:video_size" json:"video_size"`
	VideoType        int64     `gorm:"column:video_type" json:"video_type"` //0:Movie 1:TV 2:Movie(Private) 3:TV(Private)
	Location         string    `gorm:"column:location" json:"location"`     //Location on the disk
	PosterPath       string    `gorm:"column:poster_path" json:"poster_path"`
	BackdropPath     string    `gorm:"column:backdrop_path" json:"backdrop_path"`
	UserAccount      string    `gorm:"column:user_account" json:"user_account"`
	Tags             string    `gorm:"column:tags" json:"tags"`
	LikeCount        int64     `gorm:"column:like_account" json:"like_count"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time" json:"update_time"`
}

func (m *Video) Table() string {
	return "t_video"
}
