package database

import (
	"github.com/imdario/mergo"
	"github.com/pmdcosta/mmm"
)

// ensure EpisodeService implements mmm.EpisodeService.
var _ mmm.EpisodeService = &EpisodeService{}

// EpisodeService episode management service to interact with the database.
type EpisodeService struct {
	client *Client
}

// ListEpisodes returns a list of episodes.
func (s *EpisodeService) ListEpisodes(season mmm.SeasonID) ([]mmm.Episode, error) {
	// retrieve episodes.
	var v []mmm.Episode
	if err := s.client.db.Find("season", season, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return v, nil
}

// CreateEpisode persists an episode to the database.
func (s *EpisodeService) CreateEpisode(v *mmm.Episode) (*mmm.Episode, error) {
	// require object and id.
	if v == nil || v.Title == "" || v.Series == mmm.SeriesID(0) || v.Season == mmm.SeasonID(0) || v.Index == 0 {
		return nil, mmm.ErrEpisodeRequired
	}

	// update age.
	v.ModTime = s.client.Now()

	// save record.
	if err := s.client.db.Save(&v); err != nil {
		s.client.logger.Log("err", ErrDatabaseInsert, "msg", err.Error())
		return nil, ErrDatabaseInsert
	}
	return v, nil
}

// Episode retrieves an episode from the database.
func (s *EpisodeService) Episode(id mmm.EpisodeID) (*mmm.Episode, error) {
	// retrieve episode.
	var v mmm.Episode
	if err := s.client.db.One("id", id, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return &v, nil
}

// UpdateEpisode updates an episode from the db.
func (s *EpisodeService) UpdateEpisode(id mmm.EpisodeID, new *mmm.Episode) (*mmm.Episode, error) {
	// retrieve episode.
	var old mmm.Episode
	if err := s.client.db.One("id", id, &old); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}

	// merge episode with old one.
	if err := mergo.Merge(&new, old); err != nil {
		s.client.logger.Log("err", ErrDatabaseMerge, "msg", err.Error())
		return nil, ErrDatabaseMerge
	}

	// update episode.
	if err := s.client.db.Update(&new); err != nil {
		s.client.logger.Log("err", ErrDatabaseUpdate, "msg", err.Error())
		return nil, ErrDatabaseUpdate
	}
	return new, nil
}

// DeleteEpisode deletes an episode from the db.
func (s *EpisodeService) DeleteEpisode(id mmm.EpisodeID) error {
	// retrieve episode.
	var v mmm.Episode
	if err := s.client.db.One("id", id, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return ErrDatabaseQuery
	}

	// delete episode.
	if err := s.client.db.DeleteStruct(&v); err != nil {
		s.client.logger.Log("err", ErrDatabaseDelete, "msg", err.Error())
		return ErrDatabaseDelete
	}
	return nil
}
