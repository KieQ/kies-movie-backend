package download

import "github.com/anacrolix/torrent"

var downloader *torrent.Client

func init() {
	config := torrent.NewDefaultClientConfig()
	config.DataDir = "./media"
	var err error
	downloader, err = torrent.NewClient(config)
	if err != nil {
		panic(err)
	}
}

func GetDownloader() *torrent.Client {
	return downloader
}
