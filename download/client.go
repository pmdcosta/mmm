package download

import (
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/mmm"
	"sync"
	"time"
)

// Client represents a client to the download functionality.
type Client struct {
	// filename to the download target directory.
	Path string

	// timer rate in minutes.
	Rate time.Duration

	// Returns the current time.
	Now func() time.Time

	// logger
	logger log.Logger

	// Injected Services
	seasonService  mmm.SeasonService
	crawlerService mmm.CrawlerService

	// Internal Services
	downloadService DownloadService

	// routine management
	quit chan bool
	wg   sync.WaitGroup
}

// NewClient creates a new download client.
func NewClient(log log.Logger, ss mmm.SeasonService, cs mmm.CrawlerService, path string, rate time.Duration) *Client {
	c := &Client{
		logger:         log,
		seasonService:  ss,
		crawlerService: cs,
		Path:           path,
		Rate:           rate,
		Now:            time.Now,
	}
	c.downloadService.client = c
	c.quit = make(chan bool)
	return c
}

// Open opens and initializes the download scheduler.
func (c *Client) Open() error {
	c.wg.Add(1)
	go func() {
		ticker := time.NewTicker(c.Rate * time.Minute)
		defer c.wg.Done()
		for {
			select {
			case <-ticker.C:
				//c.run()
			case <-c.quit:
				ticker.Stop()
				return
			default:
				c.run()
				<-c.quit
				return
			}
		}
	}()
	return nil
}

// Close closes then underlying BoltDB database.
func (c *Client) Close() error {
	c.quit <- true
	c.wg.Wait()
	return nil
}

// run runs the download operation periodically.
func (c *Client) run() error {
	// get all the active Seasons from DB.
	seasons, err := c.seasonService.ListStateSeasons(mmm.StateRunning)
	if err != nil {
		return err
	}
	for _, s := range seasons {
		go c.downloadService.UpdateSeason(s)
	}
	return nil
}
