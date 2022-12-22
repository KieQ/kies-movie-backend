package dto

type ListDownloadStatus int64

const (
	ListDownloadStatusCannotDownload ListDownloadStatus = iota
	ListDownloadStatusCanDownload
	ListDownloadStatusDownloading
	ListDownloadStatusFinishDownload
)

type ListVideoType int64

const (
	ListVideoTypeMovie ListVideoType = iota
	ListVideoTypeTV
)
