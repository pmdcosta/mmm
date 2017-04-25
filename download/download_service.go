package download

import (
	"github.com/pmdcosta/mmm"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadService downloads new torrents for the running seasons.
type DownloadService struct {
	client *Client
}

// UpdateSeason crawls the most recent torrents for a season.
func (s *DownloadService) UpdateSeason(season mmm.Season) error {
	// crawl season website.
	torrents, err := s.client.crawlerService.Search(&season)
	if err != nil {
		return err
	}

	// check for updates.
	if len(torrents) == 0 {
		return nil
	}

	// download torrents.
	for _, t := range torrents {
		s.client.logger.Log("msg", "downloading torrent", "torrent", t.Title)
		s.downloadTorrent(t)
	}

	// save new head torrent.
	season.HeadTorrent = torrents[0].Url.String()
	s.client.seasonService.UpdateSeason(&season)

	return nil
}

// downloadTorrent downloads a torrent file.
func (s *DownloadService) downloadTorrent(t mmm.Torrent) error {
	// set file path.
	path := filepath.Join(s.client.Path, t.Title+".torrent")

	// create file.
	file, err := os.Create(path)
	if err != nil {
		s.client.logger.Log("err", ErrCreateFile, "msg", err.Error())
		return err
	}
	defer file.Close()

	// download file.
	response, err := http.Get(t.Url.String())
	if err != nil {
		s.client.logger.Log("err", ErrDownloadFile, "msg", err.Error())
		return err
	}
	defer response.Body.Close()

	// write to file.
	if _, err = io.Copy(file, response.Body); err != nil {
		s.client.logger.Log("err", ErrWriteFile, "msg", err.Error())
		return err
	}
	return nil
}
