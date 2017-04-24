package crawler

import (
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/mmm"
	"time"
)

// Client represents a client to the torrent crawlers.
type Client struct {
	// Returns the current time.
	Now func() time.Time

	// logger
	logger log.Logger
}

// NewClient creates a new db client.
func NewClient(log log.Logger) *Client {
	c := &Client{
		logger: log,
		Now:    time.Now,
	}
	return c
}

// Search searches for torrents.
func (c *Client) Search(s mmm.Season) ([]mmm.Torrent, error) {
	if s.Query == "" {
		c.logger.Log("err", mmm.ErrSeasonRequired)
		return nil, mmm.ErrSeasonRequired
	}
	switch s.Type {
	case mmm.AnimeType:
		return Search(c.logger, s.Query)
	default:
		c.logger.Log("err", mmm.ErrTypeNotFound, "msg", s.Type)
		return nil, mmm.ErrTypeNotFound
	}
}
