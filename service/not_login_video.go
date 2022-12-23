package service

import (
	"github.com/anacrolix/torrent/metainfo"
	"kies-movie-backend/download"
	"kies-movie-backend/dto"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
)

func TransForNotLoginVideoListDTO(videos []*table.Video, userInfo []*table.User) []*dto.NotLoginVideoListItem {
	userMap := make(map[string]*table.User)
	for _, item := range userInfo {
		userMap[item.Account] = item
	}

	items := make([]*dto.NotLoginVideoListItem, 0, len(videos))
	for _, oneVideo := range videos {
		oneItem := &dto.NotLoginVideoListItem{
			ID:          oneVideo.ID,
			Region:      oneVideo.Region,
			PosterPath:  oneVideo.PosterPath,
			Title:       oneVideo.VideoName,
			Description: oneVideo.VideoDescription,
		}

		//UserName and UserProfile
		if user, exist := userMap[oneVideo.UserAccount]; exist {
			oneItem.UserName = user.NickName
			oneItem.UserProfile = user.Profile
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
			oneItem.DownloadStatus = dto.ListDownloadStatusFinishDownload
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
					if oneVideo.Downloaded {
						oneItem.DownloadStatus = dto.ListDownloadStatusFinishDownload
					} else {
						oneItem.DownloadStatus = dto.ListDownloadStatusCanDownload
					}
				}
			}
		}

		items = append(items, oneItem)
	}
	return items
}
