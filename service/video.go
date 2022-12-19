package service

import (
	"context"
	"errors"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"kies-movie-backend/dto"
	"kies-movie-backend/model/table"
	"kies-movie-backend/utils"
	"mime/multipart"
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

func StartDownloadFromFile(ctx context.Context, fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		logs.CtxWarn(ctx, "fileHeader is nil")
		return errors.New("posted file is nil")
	}
	logs.CtxInfo(ctx, "filename=%v", fileHeader.Filename)
	file, err := fileHeader.Open()
	if err != nil {
		logs.CtxWarn(ctx, "failed to open file, err=%v", err)
		return err
	}

	mi, err := metainfo.Load(file)
	if err != nil {
		logs.CtxWarn(ctx, "failed to load file, err=%v", err)
		return err
	}

	client, err := torrent.NewClient(nil)
	if err != nil {
		logs.CtxWarn(ctx, "failed to create client, err=%v", err)
		return err
	}

	t, err := client.AddTorrent(mi)
	if err != nil {
		logs.CtxWarn(ctx, "failed to add torrent to client, err=%v", err)
		return err
	}

	var total int64 = 0
	for _, item := range t.Files() {
		item.Download()
		break
	}

	logs.CtxInfo(ctx, "%v, %v", t.Length(), total)
	client.WaitAll()
	return nil
}

func StartDownloadFromLink(ctx context.Context, link string) error {
	client, err := torrent.NewClient(nil)
	if err != nil {
		logs.CtxWarn(ctx, "failed to create client, err=%v", err)
		return err
	}

	t, err := client.AddMagnet(link)
	if err != nil {
		logs.CtxWarn(ctx, "failed to add torrent to client, err=%v", err)
		return err
	}

	t.DownloadAll()

	return nil
}
