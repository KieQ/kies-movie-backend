package table

import "time"

const NameVideo = "t_video"

type Video struct {
	ID               int64     `gorm:"column:id" json:"id"`
	VideoName        string    `gorm:"column:video_name" json:"video_name"`
	VideoDescription string    `gorm:"column:video_description" json:"video_description"`
	VideoType        VideoType `gorm:"column:video_type" json:"video_type"`
	Region           string    `gorm:"column:region" json:"region"`
	Link             string    `gorm:"column:link" json:"link"`
	LinkType         LinkType  `gorm:"column:link_type" json:"link_type"`
	Files            string    `gorm:"column:files" json:"files"`
	Downloaded       bool      `gorm:"column:downloaded" json:"downloaded"`
	PosterPath       string    `gorm:"column:poster_path" json:"poster_path"`
	BackdropPath     string    `gorm:"column:backdrop_path" json:"backdrop_path"`
	UserAccount      string    `gorm:"column:user_account" json:"user_account"`
	Tags             string    `gorm:"column:tags" json:"tags"`
	Liked            bool      `gorm:"column:liked" json:"liked"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time" json:"update_time"`
}

type VideoType int64

const (
	VideoTypeMovie VideoType = iota
	VideoTypeTV
	VideoTypeMoviePrivate
	VideoTypeTVPrivate
)

type LinkType int64

const (
	LinkTypeNoLink LinkType = iota
	LinkTypeMagnet
	LinkTypeLinkAddress
)
