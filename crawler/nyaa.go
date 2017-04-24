package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
	"github.com/pmdcosta/mmm"
	"net/url"
	"strconv"
)

const baseUrl = "http://www.nyaa.se/?page=search&cats=1_0&filter=0"

// Search searches torrents in nyaa.se.
func Search(logger log.Logger, query string) ([]mmm.Torrent, error) {
	// build url.
	u, err := url.Parse(baseUrl)
	if err != nil {
		logger.Log("err", ErrCrawlerFailed, "msg", err.Error())
		return nil, err
	}

	// build search query.
	q := u.Query()
	q.Set("term", query)
	u.RawQuery = q.Encode()

	// load the search results.
	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		logger.Log("err", ErrCrawlerPage, "msg", err.Error())
		return nil, err
	}

	// build torrent list.
	torrents := make([]mmm.Torrent, 0)
	doc.Find("#main > div.content > table.tlist > tbody > tr").Each(func(i int, s *goquery.Selection) {
		t, err := parseTorrent(s)
		if err != nil {
			return
		}
		torrents = append(torrents, *t)
	})
	return torrents, nil
}

func parseTorrent(s *goquery.Selection) (*mmm.Torrent, error) {
	var t mmm.Torrent
	var err error
	// get torrent url.
	d, exists := s.Find(".tlistdownload a").Attr("href")
	t.Url, err = url.Parse(d)
	if err != nil || !exists {
		return &t, err
	}
	t.Url.Scheme = "https"
	// get torrent seeds.
	if t.Seeds, err = strconv.Atoi(s.Find(".tlistsn").Text()); err != nil {
		return &t, err
	}
	// get torrent leeches.
	if t.Leeches, err = strconv.Atoi(s.Find(".tlistln").Text()); err != nil {
		return &t, err
	}
	// get torrent title.
	t.Title = s.Find(".tlistname a").Text()
	// get torrent size.
	t.Size = s.Find(".tlistsize").Text()
	return &t, nil
}
