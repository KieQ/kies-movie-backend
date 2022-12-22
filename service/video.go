package service

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"kies-movie-backend/download"
	"kies-movie-backend/dto"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
)

func TransForVideoListDTO(videos []*table.Video) []*dto.VideoListItem {
	items := make([]*dto.VideoListItem, 0, len(videos))
	for _, oneVideo := range videos {
		oneItem := &dto.VideoListItem{
			ID:          oneVideo.ID,
			Region:      oneVideo.Region,
			PosterPath:  oneVideo.PosterPath,
			Title:       oneVideo.VideoName,
			Liked:       oneVideo.Liked,
			Description: oneVideo.VideoDescription,
		}
		// IsPublic
		if utils.Contain([]table.VideoType{table.VideoTypeMovie, table.VideoTypeTV}, oneVideo.VideoType) {
			oneItem.IsPublic = true
		} else {
			oneItem.IsPublic = false
		}

		//VideoType
		if utils.Contain([]table.VideoType{table.VideoTypeMovie, table.VideoTypeMoviePrivate}, oneVideo.VideoType) {
			oneItem.VideoType = dto.ListVideoTypeMovie
		} else if utils.Contain([]table.VideoType{table.VideoTypeTV, table.VideoTypeTVPrivate}, oneVideo.VideoType) {
			oneItem.VideoType = dto.ListVideoTypeTV
		}

		//DownloadStatus
		if oneVideo.LinkType == table.LinkTypeNoLink {
			oneItem.DownloadStatus = dto.ListDownloadStatusCannotDownload
		} else if oneVideo.LinkType == table.LinkTypeLinkAddress {
			oneItem.DownloadStatus = dto.ListDownloadStatusCannotDownload
		} else if oneVideo.LinkType == table.LinkTypeMagnet {
			mag, err := metainfo.ParseMagnetUri(oneVideo.Link)
			if err != nil {
				oneItem.DownloadStatus = dto.ListDownloadStatusCannotDownload
			} else {
				t, exist, err := download.GetFromDownloadingMap(mag.InfoHash.HexString())
				if err != nil {
					oneItem.DownloadStatus = dto.ListDownloadStatusCannotDownload
				} else if exist && len(t.DownloadingFiles) != 0 {
					if t.AllFinished() {
						oneItem.DownloadStatus = dto.ListDownloadStatusFinishDownload
					} else if t.AllPause() {
						oneItem.DownloadStatus = dto.ListDownloadStatusCanDownload
					} else {
						oneItem.DownloadStatus = dto.ListDownloadStatusDownloading
					}
				} else {
					oneItem.DownloadStatus = dto.ListDownloadStatusCanDownload
				}
			}
		}

		//CanPlayFiles
		if oneVideo.LinkType == table.LinkTypeLinkAddress {
			oneItem.CanPlayFiles = []*dto.CanPlayFilesItem{{
				Path:            oneVideo.Link,
				DisplayPath:     oneVideo.VideoName,
				DownloadedBytes: 0,
				TotalBytes:      0,
				CanPlay:         true,
				Downloading:     false,
			}}
		} else if oneVideo.LinkType == table.LinkTypeMagnet {
			var useDB bool
			mag, err := metainfo.ParseMagnetUri(oneVideo.Link)
			if err != nil {
				useDB = true
			} else {
				t, exist, err := download.GetFromDownloadingMap(mag.InfoHash.HexString())
				if err == nil && exist {
					for _, file := range t.DownloadingFiles {
						oneItem.CanPlayFiles = append(oneItem.CanPlayFiles, &dto.CanPlayFilesItem{
							Path:            file.Path(),
							DisplayPath:     file.DisplayPath(),
							DownloadedBytes: file.BytesCompleted(),
							TotalBytes:      file.Length(),
							CanPlay:         file.BytesCompleted() == file.Length(),
							Downloading:     file.Priority() != torrent.PiecePriorityNone,
						})
					}
				} else {
					useDB = true
				}
			}
			if useDB && oneVideo.Files != "" && oneVideo.Downloaded {
				files := utils.FromJSON[[]string](oneVideo.Files)
				for _, file := range files {
					fileSize := download.FileSize(file)
					oneItem.CanPlayFiles = []*dto.CanPlayFilesItem{{
						Path:            file,
						DisplayPath:     download.GetNaiveDisplayPath(file),
						DownloadedBytes: fileSize,
						TotalBytes:      fileSize,
						CanPlay:         fileSize != 0,
						Downloading:     false,
					}}
				}
			}

		}

		items = append(items, oneItem)
	}
	return items
}
