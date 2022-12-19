package service

import (
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
			oneItem.VideoType = dto.VideoListItemVideoMovie
		} else if utils.Contain([]table.VideoType{table.VideoTypeTV, table.VideoTypeTVPrivate}, oneVideo.VideoType) {
			oneItem.VideoType = dto.VideoListItemVideoTV
		}

		//Status
		if oneVideo.Location == "" {
			oneItem.Status = dto.VideoListItemStatusUnableToDownload
		} else {
			oneItem.Status = dto.VideoListItemStatusCanDownload
			//TODO Downloader logic
		}

		items = append(items, oneItem)
	}
	return items
}
