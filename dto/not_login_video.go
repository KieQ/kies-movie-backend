package dto

type NotLoginVideoListItem struct {
	ID             int64              `json:"id"`
	DownloadStatus ListDownloadStatus `json:"download_status"`
	Region         string             `json:"region"`
	PosterPath     string             `json:"poster_path"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	VideoType      ListVideoType      `json:"video_type"`
	UserName       string             `json:"user_name"`
	UserProfile    string             `json:"user_profile"`
}

type NotLoginVideoListResponse struct {
	Total int64                    `json:"total"`
	Page  int32                    `json:"page"`
	Size  int32                    `json:"size"`
	Items []*NotLoginVideoListItem `json:"items"`
}
