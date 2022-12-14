package dto

type SessionLoginRequest struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type SessionLoginResponse struct {
	NickName string `json:"nick_name"`
}

type SessionSignupRequest struct {
	Account          string   `json:"account"`
	Password         string   `json:"password"`
	NickName         string   `json:"nick_name"`
	Gender           int32    `json:"gender"`
	Profile          string   `json:"profile"`
	Phone            string   `json:"phone"`
	Email            string   `json:"email"`
	SelfIntroduction string   `json:"self_introduction"`
	PreferTags       []string `json:"prefer_tags"`
	DefaultLanguage  string   `json:"default_language"`
}

type SessionLogoutRequest struct {
	Account string `json:"account"`
}
