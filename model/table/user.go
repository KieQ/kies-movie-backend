package table

import "time"

const NameUser = "t_user"

type User struct {
	ID               int64     `gorm:"column:id" json:"id"`
	Account          string    `gorm:"column:account" json:"account"`
	Password         string    `gorm:"column:password" json:"password"`
	NickName         string    `gorm:"column:nick_name" json:"nick_name"`
	Profile          string    `gorm:"column:profile" json:"profile"`
	Phone            string    `gorm:"column:phone" json:"phone"`
	Email            string    `gorm:"column:email" json:"email"`
	Gender           Gender    `gorm:"column:gender" json:"gender"`
	SelfIntroduction string    `gorm:"column:self_introduction" json:"self_introduction"`
	PreferTags       string    `gorm:"column:prefer_tags" json:"prefer_tags"`
	DefaultLanguage  string    `gorm:"column:default_language" json:"default_language"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time" json:"update_time"`
}

type Gender int32

const (
	Male Gender = iota
	Female
)
