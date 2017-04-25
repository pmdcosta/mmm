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

	// services
	crawlerService CrawlerService

	// logger
	logger log.Logger
}

// NewClient creates a new db client.
func NewClient(log log.Logger) *Client {
	c := &Client{
		logger: log,
		Now:    time.Now,
	}
	c.crawlerService.client = c
	c.crawlerService.nyaa.client = c
	return c
}

// CrawlerService returns the created crawler service associated with the client.
func (c *Client) CrawlerService() mmm.CrawlerService { return &c.crawlerService }
