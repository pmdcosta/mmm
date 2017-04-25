package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pmdcosta/mmm"
	"net/url"
	"strconv"
)

// Nyaa represents a crawler source at nyaa.se.
type Nyaa struct {
	client *Client
}

const NyaaUrl = "http://www.nyaa.se/?page=search&cats=1_0&filter=0"

// Search searches torrents in nyaa.se.
func (s *Nyaa) Search(season *mmm.Season) ([]mmm.Torrent, error) {
	// build url.
	u, err := url.Parse(NyaaUrl)
	if err != nil {
		s.client.logger.Log("err", ErrCrawlerFailed, "msg", err.Error())
		return nil, err
	}

	// build search query.
	q := u.Query()
	q.Set("term", season.Query)
	u.RawQuery = q.Encode()

	// load the search results.
	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		s.client.logger.Log("err", ErrCrawlerPage, "msg", err.Error())
		return nil, err
	}

	// build torrent list.
	torrents := make([]mmm.Torrent, 0)
	doc.Find("#main > div.content > table.tlist > tbody > tr").EachWithBreak(func(i int, q *goquery.Selection) bool {
		t, err := s.parseTorrent(q)
		if err != nil {
			s.client.logger.Log("err", ErrCrawlerField, "msg", err.Error())
			return true
		}
		if t == nil {
			return true
		}
		if season.HeadTorrent == t.Url.String() {
			return false
		}
		torrents = append(torrents, *t)
		return true
	})
	return torrents, nil
}

func (s *Nyaa) parseTorrent(q *goquery.Selection) (*mmm.Torrent, error) {
	var t mmm.Torrent
	var err error

	// get torrent url.
	d, exists := q.Find(".tlistdownload a").Attr("href")
	t.Url, err = url.Parse(d)
	if err != nil || !exists {
		return nil, err
	}
	t.Url.Scheme = "https"

	// get torrent seeds.
	if t.Seeds, err = strconv.Atoi(q.Find(".tlistsn").Text()); err != nil {
		return nil, err
	}

	// get torrent leeches.
	if t.Leeches, err = strconv.Atoi(q.Find(".tlistln").Text()); err != nil {
		return nil, err
	}

	// get torrent title.
	t.Title = q.Find(".tlistname a").Text()

	// get torrent size.
	t.Size = q.Find(".tlistsize").Text()
	return &t, nil
}
