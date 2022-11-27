package table

import "time"

type User struct {
	ID               int64     `gorm:"column:id" json:"id"`
	Account          string    `gorm:"column:account" json:"account"`
	Password         string    `gorm:"column:password" json:"password"`
	NickName         string    `gorm:"column:nick_name" json:"nick_name"`
	Profile          string    `gorm:"column:profile" json:"profile"`
	Phone            string    `gorm:"column:phone" json:"phone"`
	Email            string    `gorm:"column:email" json:"email"`
	Gender           int32     `gorm:"column:gender" json:"gender"`
	SelfIntroduction string    `gorm:"column:self_introduction" json:"self_introduction"`
	PreferTags       string    `gorm:"column:prefer_tags" json:"prefer_tags"`
	CreateTime       time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime       time.Time `gorm:"column:update_time" json:"update_time"`
}

func (u *User) Table() string {
	return "t_user"
}
