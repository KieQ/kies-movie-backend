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
	FailedToFindMovieOrTV
	FailedToFindUsers
	FailedToAddVideo
	CannotCloneYourOwnMovie
	NoLinkCannotBeProcessed
)

var translate = map[SentenceIndex]sentence{
	NoLinkCannotBeProcessed: {
		English: "video has no link, so you can't do this",
		Chinese: "视频没有链接，因此不能进行此操作",
	},
	CannotCloneYourOwnMovie: {
		English: "you cannot clone your own movie/TV",
		Chinese: "你不能添加自己的电影/电视剧",
	},
	FailedToAddVideo: {
		English: "failed to add video to database",
		Chinese: "添加视频到数据库中失败",
	},
	FailedToFindUsers: {
		English: "service failed to fetch user info",
		Chinese: "系统获取用户信息失败",
	},
	FailedToFindMovieOrTV: {
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
