package dto

type UserUpdateRequest struct {
	Account          string   `json:"account"`
	Password         *string  `json:"password"`
	NickName         *string  `json:"nick_name"`
	Gender           *int32   `json:"gender"`
	Profile          *string  `json:"profile"`
	Phone            *string  `json:"phone"`
	Email            *string  `json:"email"`
	SelfIntroduction *string  `json:"self_introduction"`
	PreferTags       []string `json:"prefer_tags"`
}

type User struct {
	Account          string   `json:"account"`
	NickName         string   `json:"nick_name"`
	Profile          string   `json:"profile"`
	Phone            string   `json:"phone"`
	Email            string   `json:"email"`
	Gender           int32    `json:"gender"`
	SelfIntroduction string   `json:"self_introduction"`
	PreferTags       []string `json:"prefer_tags"`
	CreateTime       int64    `json:"create_time"`
}

type UserDetailResponse struct {
	User *User `json:"user"`
}

type UserListResponse struct {
	Users []*User `json:"users"`
}
