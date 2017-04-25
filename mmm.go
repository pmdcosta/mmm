package mmm

import (
	"net/url"
	"time"
)

// Client creates a connection to the services.
type Client interface {
	SeasonService() SeasonService
	EpisodeService() EpisodeService
}

// SeasonID represents a season identifier.
type SeasonID int

// Season represents a collection of episodes that correspond to a season.
type Season struct {
	ID          SeasonID    `json:"id" storm:"id,increment"`
	Title       string      `json:"title" storm:"unique"`
	Query       string      `json:"query"`
	Type        MediaType   `json:"type" storm:"index"`
	Index       int         `json:"index,omitempty"`
	State       SeasonState `json:"state" storm:"index"`
	HeadTorrent string      `json:"head"`
	Tags        []Tag       `json:"tags" storm:"index"`
	ModTime     time.Time   `json:"modTime"`
}

// SeasonService represents a service for managing season.
type SeasonService interface {
	ListSeasons() ([]Season, error)
	ListStateSeasons(state SeasonState) ([]Season, error)
	CreateSeason(season *Season) (*Season, error)
	Season(id SeasonID) (*Season, error)
	UpdateSeason(season *Season) (*Season, error)
	DeleteSeason(id SeasonID) error
}

// EpisodeID represents a episode identifier.
type EpisodeID int

// Episode represents a single episode.
type Episode struct {
	ID      EpisodeID `json:"id" storm:"id,increment"`
	Title   string    `json:"title"`
	Index   float32   `json:"index"`
	Season  SeasonID  `json:"season" storm:"index"`
	ModTime time.Time `json:"modTime"`
}

// EpisodeService represents a service for managing episode.
type EpisodeService interface {
	ListEpisodes(season SeasonID) ([]Episode, error)
	CreateEpisode(episode *Episode) (*Episode, error)
	Episode(id EpisodeID) (*Episode, error)
	UpdateEpisode(episode *Episode) (*Episode, error)
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
	TypeAnime  MediaType = "anime"
	TypeMovie            = "movie"
	TypeSeries           = "series"
)

// SeasonState represents all the possible states for a season.
type SeasonState string

const (
	StateRunning  SeasonState = "running"
	StateComplete             = "complete"
	StateHiatus               = "hiatus"
)

// Torrent represents a downloadable torrent file.
type Torrent struct {
	Url     *url.URL `json:"url" storm:"id"`
	Title   string   `json:"title"`
	Size    string   `json:"size"`
	Seeds   int      `json:"seeds"`
	Leeches int      `json:"leeches"`
}

// CrawlerService represents a service for crawling websites.
type CrawlerService interface {
	Search(season *Season) ([]Torrent, error)
}
