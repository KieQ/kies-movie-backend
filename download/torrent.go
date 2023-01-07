package download

import (
	"context"
	"errors"
	"github.com/Kidsunbo/kie_toolbox_go/container"
	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"github.com/anacrolix/torrent"
	"kies-movie-backend/model/db"
	"kies-movie-backend/utils"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const DataDir = "./media/"

var downloader *torrent.Client

type DownloadingInfo struct {
	DownloadingFiles []*torrent.File
	Torrent          *torrent.Torrent
}

func (d *DownloadingInfo) AllFinished() bool {
	for _, file := range d.DownloadingFiles {
		if file.BytesCompleted() < file.Length() {
			return false
		}
	}
	return true
}

func (d *DownloadingInfo) AllPause() bool {
	for _, file := range d.DownloadingFiles {
		if file.Priority() != torrent.PiecePriorityNone {
			return false
		}
	}
	return true
}

var downloadingMap sync.Map

func init() {
	config := torrent.NewDefaultClientConfig()
	config.DataDir = DataDir
	config.Seed = false
	var err error
	downloader, err = torrent.NewClient(config)
	if err != nil {
		panic(err)
	}
	go startWait()
}

func startWait() {
	tick := time.Tick(5 * time.Second)
	for range tick {
		downloader.WaitAll()
	}
}

func GetFromDownloadingMap(infoHash string) (*DownloadingInfo, bool, error) {
	if v, ok := downloadingMap.Load(infoHash); ok {
		if item, yes := v.(*DownloadingInfo); yes {
			return item, true, nil
		} else {
			return nil, true, errors.New("type is not Torrent")
		}
	}
	return nil, false, nil
}

//ShowFilesInMagnet returns InfoHash, Files, timeout and error
func ShowFilesInMagnet(ctx context.Context, filesInDB, link string) (string, []*torrent.File, bool, error) {
	t, err := downloader.AddMagnet(link)
	if err != nil {
		logs.CtxInfo(ctx, "failed to add link=%v", link)
		return "", nil, false, err
	}

	v, exist, err := GetFromDownloadingMap(t.InfoHash().HexString())
	if err != nil {
		return "", nil, false, err
	} else if !exist {
		v = &DownloadingInfo{
			DownloadingFiles: nil,
			Torrent:          t,
		}
	} else {
		t = v.Torrent
	}

	timer := time.After(time.Minute)
	select {
	case <-t.GotInfo():
		if filesInDB != "" && !exist {
			files := utils.FromJSON[[]string](filesInDB)
			set := container.NewSet[string]()
			for _, item := range files {
				set.Add(item)
			}
			for _, item := range t.Files() {
				if set.Contain(item.Path()) {
					v.DownloadingFiles = append(v.DownloadingFiles, item)
				}
			}
		}
		downloadingMap.Store(t.InfoHash().HexString(), v)
	case <-timer:
		return "", nil, true, nil
	}
	return t.InfoHash().HexString(), t.Files(), false, nil
}

//StartDownloadSelectFileAsync returns downloaded files, exist in Map and error
func StartDownloadSelectFileAsync(ctx context.Context, id int64, account, infoHash string, files []string) ([]*torrent.File, bool, error) {
	//Get from downloading map
	item, exist, err := GetFromDownloadingMap(infoHash)
	if err != nil {
		return nil, false, err
	} else if !exist {
		return nil, false, nil
	}

	//Get Last time downloading files
	lastTime := make(map[string]*torrent.File)
	for _, v := range item.DownloadingFiles {
		lastTime[v.Path()] = v
	}

	//Create a set and mark the last downloading file to nil if it needs to be downloaded this time
	set := container.NewSet[string]()
	for _, file := range files {
		set.Add(file)
		if _, existInMap := lastTime[file]; existInMap {
			lastTime[file] = nil
		}
	}

	//Start download the files
	item.DownloadingFiles = nil
	for _, file := range item.Torrent.Files() {
		if set.Contain(file.Path()) {
			item.DownloadingFiles = append(item.DownloadingFiles, file)
			file.Download()
		}
	}

	//Delete the no needed files
	stoppedFile := make([]*torrent.File, 0, len(lastTime))
	for _, v := range lastTime {
		if v != nil {
			stoppedFile = append(stoppedFile, v)
		}
	}

	deleteFiles := StopDownload(ctx, stoppedFile)
	stillThere := DeleteFiles(ctx, deleteFiles)

	go func() {
		tick := time.Tick(time.Minute)
		for range tick {
			downloadOK := false
			deleteOK := false
			//Delete
			if len(stillThere) != 0 {
				stillThere = DeleteFiles(ctx, stillThere)
			} else {
				deleteOK = true
			}

			//Download
			if item.AllFinished() {
				rows, err := db.UpdateVideoByID(ctx, account, id, map[string]interface{}{
					"downloaded": true,
				})
				if err != nil || rows == 0 {
					logs.CtxWarn(ctx, "failed to update downloaded, err=%v, rows=%v", err, rows)
					continue
				}
				item.Torrent.Drop()
				downloadingMap.Delete(infoHash)
				downloadOK = true
			}

			//Check if flow need return
			if deleteOK && downloadOK {
				return
			}
		}
	}()

	return item.DownloadingFiles, true, nil
}

func StopDownload(ctx context.Context, files []*torrent.File) []string {
	stoppedFile := make([]string, 0, len(files))
	for _, item := range files {
		item.SetPriority(torrent.PiecePriorityNone)
		fileLocation := WrapPath(item.Path())
		stoppedFile = append(stoppedFile, fileLocation)
	}
	logs.CtxInfo(ctx, "stopped files: %v", utils.ToJSON(stoppedFile))
	return stoppedFile
}

func WrapPath(filePath string) string {
	return path.Join(DataDir, filePath)
}

//DeleteFiles delete the files and return the undeleted files，be sure the path is exact one.
func DeleteFiles(ctx context.Context, files []string) []string {
	stillThere := make([]string, 0, len(files))
	for _, fileLocation := range files {
		err := os.Remove(fileLocation)
		if err != nil {
			logs.CtxError(ctx, "failed to delete file, err=%v", err)
		}
		_, e := os.Stat(fileLocation)
		if !os.IsNotExist(e) {
			stillThere = append(stillThere, fileLocation)
		}
	}
	logs.CtxInfo(ctx, "still there files: %v", utils.ToJSON(stillThere))
	return stillThere
}

//DeleteWholeDirectory delete the whole directory and return the undeleted files and directories.
func DeleteWholeDirectory(ctx context.Context, files []string) []string {
	set := container.NewSet[string]()
	for _, file := range files {
		file = strings.TrimLeft(file, "/")
		paths := strings.SplitN(file, "/", 2)
		if len(paths) == 0 {
			continue
		} else {
			set.Add(paths[0])
		}
	}

	stillThere := make([]string, 0, set.Size())
	for fileLocation := range set.Range() {
		err := os.RemoveAll(fileLocation)
		if err != nil {
			logs.CtxError(ctx, "failed to delete file, err=%v", err)
		}
		_, e := os.Stat(fileLocation)
		if !os.IsNotExist(e) {
			stillThere = append(stillThere, fileLocation)
		}
	}
	return stillThere
}

func GetNaiveDisplayPath(path string) string {
	result := strings.SplitN(path, "/", 2)
	if len(result) == 2 {
		return result[1]
	}
	return path
}

func FileSize(nativePath string) int64 {
	file, err := os.Stat(WrapPath(nativePath))
	if err != nil {
		return 0
	}
	return file.Size()
}
