package i18n

import "kies-movie-backend/constant"

type sentence struct {
	English string
	Chinese string
}

type SentenceIndex int32

const (
	UserHasExisted SentenceIndex = iota
	FailedToCheckExistence
	FailedToLogin
	CouldNotFindUser
	NoMovieOrTVFound
)

var translate = map[SentenceIndex]sentence{
	NoMovieOrTVFound: {
		English: "service failed to fetch movies and tv",
		Chinese: "系统获取电影/电视剧信息失败",
	},
	CouldNotFindUser: {
		English: "service could not find user information, please check",
		Chinese: "系统无法找到用户信息，请检查",
	},
	FailedToLogin: {
		English: "failed to login, username or password is wrong",
		Chinese: "登陆失败，用户名或者密码错误",
	},
	FailedToCheckExistence: {
		English: "server failed to check the existence of user, please try again",
		Chinese: "系统无法检查用户是否存在，请重试",
	},
	UserHasExisted: {
		English: "service can't perform this signup because user has existed",
		Chinese: "由于用户已存在，系统无法完成注册",
	},

	SentenceIndex(constant.ServiceError): {
		English: "service error",
		Chinese: "服务器错误",
	},
	SentenceIndex(constant.UserNotLogin): {
		English: "user has not logged in",
		Chinese: "用户未登陆",
	},
	SentenceIndex(constant.UserIPChanged): {
		English: "user's ip has changed, please login again",
		Chinese: "用户IP改变，请重新登录",
	},
	SentenceIndex(constant.RequestParameterError): {
		English: "the request parameter is illegal",
		Chinese: "请求参数不合法",
	},
	SentenceIndex(constant.FailedToProcess): {
		English: "service failed to process this request",
		Chinese: "本次请求处理失败",
	},
	SentenceIndex(constant.NoAuthority): {
		English: "user has no authority to operate",
		Chinese: "用户无权操作",
	},
}
