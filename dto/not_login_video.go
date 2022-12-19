package dto

type NotLoginVideoListItemStatus int64

const (
	NotLoginVideoListItemStatusUnableToDownload VideoListItemStatus = iota
	NotLoginVideoListItemStatusCanDownload
	NotLoginVideoListItemStatusDownloading
	NotLoginVideoListItemStatusCanPlay
)

type NotLoginVideoListItemVideoType int64

const (
	NotLoginVideoListItemVideoMovie VideoListItemVideoType = iota
	NotLoginVideoListItemVideoTV
)

type NotLoginVideoListItem struct {
	ID          int64                  `json:"id"`
	Status      VideoListItemStatus    `json:"status"`
	Region      string                 `json:"region"`
	PosterPath  string                 `json:"poster_path"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	VideoType   VideoListItemVideoType `json:"video_type"`
	UserName    string                 `json:"user_name"`
	UserProfile string                 `json:"user_profile"`
}

type NotLoginVideoListResponse struct {
	Total int64                    `json:"total"`
	Page  int32                    `json:"page"`
	Size  int32                    `json:"size"`
	Items []*NotLoginVideoListItem `json:"items"`
}
