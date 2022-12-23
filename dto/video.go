package dto

type CanPlayFilesItem struct {
	Path             string `json:"path"`
	DisplayPath      string `json:"display_path"`
	DownloadedBytes  int64  `json:"downloaded_bytes"`
	TotalBytes       int64  `json:"total_bytes"`
	CanPlay          bool   `json:"can_play"`
	Downloading      bool   `json:"downloading"`
	VideoOnOtherSite bool   `json:"video_on_other_site"`
}

type VideoListItem struct {
	ID             int64               `json:"id"`
	DownloadStatus ListDownloadStatus  `json:"download_status"`
	CanPlayFiles   []*CanPlayFilesItem `json:"can_play_files"`
	IsPublic       bool                `json:"is_public"`
	Region         string              `json:"region"`
	PosterPath     string              `json:"poster_path"`
	Title          string              `json:"title"`
	Liked          bool                `json:"liked"`
	Description    string              `json:"description"`
	VideoType      ListVideoType       `json:"video_type"`
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

type VideoAvailableFileInfo struct {
	Path            string `json:"path"`
	DisplayPath     string `json:"display_path"`
	DownloadedBytes int64  `json:"downloaded_bytes"`
	TotalBytes      int64  `json:"total_bytes"`
	Downloading     bool   `json:"downloading"`
}

type VideoAvailableFilesResponse struct {
	Timeout  bool                     `json:"timeout"`
	InfoHash string                   `json:"info_hash"`
	Files    []VideoAvailableFileInfo `json:"files"`
}

type VideoDownloadRequest struct {
	ID       int64    `json:"id"`
	InfoHash string   `json:"info_hash"`
	Files    []string `json:"files"`
}
