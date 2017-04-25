package crawler

import "github.com/pmdcosta/mmm"

// CrawlerService crawls websites for torrents.
type CrawlerService struct {
	client *Client

	// Crawler sources
	nyaa Nyaa
}

// Search searches for torrents.
func (s *CrawlerService) Search(season *mmm.Season) ([]mmm.Torrent, error) {
	if season.Query == "" {
		s.client.logger.Log("err", mmm.ErrSeasonRequired)
		return nil, mmm.ErrSeasonRequired
	}

	switch season.Type {
	case mmm.TypeAnime:
		return s.nyaa.Search(season)
	default:
		s.client.logger.Log("err", mmm.ErrTypeNotFound, "msg", season.Type)
		return nil, mmm.ErrTypeNotFound
	}
}
