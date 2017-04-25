package crawler

import "github.com/pmdcosta/mmm"

// Crawler errors.
const (
	ErrCrawlerFailed = mmm.Error("failed to start crawler")
	ErrCrawlerPage   = mmm.Error("failed to crawl web page")
	ErrCrawlerField  = mmm.Error("failed to crawl field")
)
