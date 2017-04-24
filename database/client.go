package database

import (
	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/mmm"
	"time"
)

// Client represents a client to the underlying BoltDB data store.
type Client struct {
	// Filename to the BoltDB database.
	Path string

	// Returns the current time.
	Now func() time.Time

	// Services
	seriesService  SeriesService
	seasonService  SeasonService
	episodeService EpisodeService

	// logger
	logger log.Logger

	// Bolt DB
	db *storm.DB
}

// NewClient creates a new db client.
func NewClient(log log.Logger, path string) *Client {
	c := &Client{
		logger: log,
		Path:   path,
		Now:    time.Now,
	}
	c.seriesService.client = c
	c.seasonService.client = c
	c.episodeService.client = c
	return c
}

// Open opens and initializes the BoltDB database.
func (c *Client) Open() error {
	// Open database file.
	db, err := storm.Open(c.Path, storm.BoltOptions(0600, &bolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		c.logger.Log("err", ErrDatabaseFailed, "msg", err.Error())
		return err
	}
	c.db = db
	return nil
}

// Close closes then underlying BoltDB database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// SeriesService returns the series service associated with the client.
func (c *Client) SeriesService() mmm.SeriesService { return &c.seriesService }

// SeasonService returns the season service associated with the client.
func (c *Client) SeasonService() mmm.SeasonService { return &c.seasonService }

// EpisodeService returns the episode service associated with the client.
func (c *Client) EpisodeService() mmm.EpisodeService { return &c.episodeService }
