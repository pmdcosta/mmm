package mmm

import (
	"net/url"
	"time"
)

// Client creates a connection to the services.
type Client interface {
	SeriesService() SeriesService
	SeasonService() SeasonService
	EpisodeService() EpisodeService
}

// SeriesID represents a series identifier.
type SeriesID int

// Series represents a collection of episodes that correspond to a series.
type Series struct {
	ID      SeriesID  `json:"id" storm:"id,increment"`
	Title   string    `json:"title" storm:"index"`
	Type    MediaType `json:"type" storm:"index"`
	Tags    []Tag     `json:"tags" storm:"index"`
	ModTime time.Time `json:"-"`
}

// SeriesService represents a service for managing series.
type SeriesService interface {
	ListSeries() ([]Series, error)
	CreateSeries(series *Series) (*Series, error)
	Series(id SeriesID) (*Series, error)
	UpdateSeries(id SeriesID, series *Series) (*Series, error)
	DeleteSeries(id SeriesID) error
}

// SeasonID represents a season identifier.
type SeasonID int

// Season represents a collection of episodes that correspond to a season.
type Season struct {
	ID       SeasonID  `json:"id" storm:"id,increment"`
	Title    string    `json:"title" storm:"index"`
	Query    string    `json:"query"`
	Type     MediaType `json:"type" storm:"index"`
	Index    int       `json:"index,omitempty"`
	Series   SeriesID  `json:"series" storm:"index"`
	Complete bool      `json:"complete" storm:"index"`
	ModTime  time.Time `json:"-"`
}

// SeasonService represents a service for managing season.
type SeasonService interface {
	ListSeasons(series SeriesID) ([]Season, error)
	ListCompleteSeasons(complete bool) ([]Season, error)
	CreateSeason(season *Season) (*Season, error)
	Season(id SeasonID) (*Season, error)
	UpdateSeason(id SeasonID, season *Season) (*Season, error)
	DeleteSeason(id SeasonID) error
}

// EpisodeID represents a episode identifier.
type EpisodeID int

// Episode represents a single episode.
type Episode struct {
	ID      EpisodeID `json:"id" storm:"id,increment"`
	Title   string    `json:"title"`
	Index   float32   `json:"index"`
	Series  SeriesID  `json:"series" storm:"index"`
	Season  SeasonID  `json:"season" storm:"index"`
	ModTime time.Time `json:"-"`
}

// EpisodeService represents a service for managing episode.
type EpisodeService interface {
	ListEpisodes(season SeasonID) ([]Episode, error)
	CreateEpisode(episode *Episode) (*Episode, error)
	Episode(id EpisodeID) (*Episode, error)
	UpdateEpisode(id EpisodeID, episode *Episode) (*Episode, error)
	DeleteEpisode(id EpisodeID) error
}

// Tag represents a common type of media content.
type Tag struct {
	ID string
}

// Tag represents a service for managing media types.
type TagService interface {
	ListTags() ([]Tag, error)
	CreateTag(tag *Tag) error
	DeleteTag(tag *Tag) error
}

// MediaType represents all the possible media types.
type MediaType string

const (
	AnimeType  MediaType = "anime"
	MovieType            = "movie"
	SeriesType           = "series"
)

// Torrent represents a downloadable torrent file.
type Torrent struct {
	Url     *url.URL `json:"url" storm:"id"`
	Title   string   `json:"title"`
	Size    string   `json:"size"`
	Seeds   int      `json:"seeds"`
	Leeches int      `json:"leeches"`
}
