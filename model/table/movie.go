package table

import "time"

type Movie struct {
	ID         int64     `gorm:"column:id" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	MagnetLink string    `gorm:"column:magnet_link" json:"magnet_link"`
	Size       int64     `gorm:"column:size" json:"size"`
	Location   string    `gorm:"column:location" json:"location"`
	Profile    string    `gorm:"column:profile" json:"profile"`
	UserID     string    `gorm:"column:user_id" json:"user_id"`
	IsPrivate  int32     `gorm:"column:is_private" json:"is_private"`
	Tags       string    `gorm:"column:tags" json:"tags"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (m *Movie) Table() string {
	return "t_movie"
}
