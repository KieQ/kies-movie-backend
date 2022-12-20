package dto

type VideoListItemStatus int64

const (
	VideoListItemStatusUnableToDownload VideoListItemStatus = iota
	VideoListItemStatusCanDownload
	VideoListItemStatusDownloading
	VideoListItemStatusCanPlay
)

type VideoListItemVideoType int64

const (
	VideoListItemVideoMovie VideoListItemVideoType = iota
	VideoListItemVideoTV
)

type VideoListItem struct {
	ID          int64                  `json:"id"`
	Status      VideoListItemStatus    `json:"status"`
	IsPublic    bool                   `json:"is_public"`
	Region      string                 `json:"region"`
	PosterPath  string                 `json:"poster_path"`
	Title       string                 `json:"title"`
	Liked       bool                   `json:"liked"`
	Description string                 `json:"description"`
	VideoType   VideoListItemVideoType `json:"video_type"`
}

type VideoListResponse struct {
	Total int64            `json:"total"`
	Page  int32            `json:"page"`
	Size  int32            `json:"size"`
	Items []*VideoListItem `json:"items"`
}

type VideoDeleteRequest struct {
	ID int64 `json:"id"`
}

type VideoLikeRequest struct {
	ID    int64 `json:"id"`
	Liked bool  `json:"liked"`
}

type VideoCloneRequest struct {
	ID int64 `json:"id"`
}
