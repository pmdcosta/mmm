package database

import (
	"github.com/imdario/mergo"
	"github.com/pmdcosta/mmm"
)

// ensure SeasonService implements mmm.SeasonService.
var _ mmm.SeasonService = &SeasonService{}

// SeasonService season management service to interact with the database.
type SeasonService struct {
	client *Client
}

// ListSeasons returns a list of seasons.
func (s *SeasonService) ListSeasons(series mmm.SeriesID) ([]mmm.Season, error) {
	// retrieve seasons.
	var v []mmm.Season
	if err := s.client.db.Find("series", series, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return v, nil
}

// ListCompleteSeasons returns a list of complete/incomplete seasons.
func (s *SeasonService) ListCompleteSeasons(b bool) ([]mmm.Season, error) {
	// retrieve seasons.
	var v []mmm.Season
	if err := s.client.db.Find("complete", b, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return v, nil
}

// CreateSeason persists a season to the database.
func (s *SeasonService) CreateSeason(v *mmm.Season) (*mmm.Season, error) {
	// require object and id.
	if v == nil || v.Title == "" || v.Series == mmm.SeriesID(0) || v.Index == 0 {
		return nil, mmm.ErrSeasonRequired
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

// Season retrieves a season from the database.
func (s *SeasonService) Season(id mmm.SeasonID) (*mmm.Season, error) {
	// retrieve season.
	var v mmm.Season
	if err := s.client.db.One("id", id, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return &v, nil
}

// UpdateSeason updates a season from the db.
func (s *SeasonService) UpdateSeason(id mmm.SeasonID, new *mmm.Season) (*mmm.Season, error) {
	// retrieve season.
	var old mmm.Season
	if err := s.client.db.One("id", id, &old); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}

	// merge season with old one.
	if err := mergo.Merge(&new, old); err != nil {
		s.client.logger.Log("err", ErrDatabaseMerge, "msg", err.Error())
		return nil, ErrDatabaseMerge
	}

	// update season.
	if err := s.client.db.Update(&new); err != nil {
		s.client.logger.Log("err", ErrDatabaseUpdate, "msg", err.Error())
		return nil, ErrDatabaseUpdate
	}
	return new, nil
}

// DeleteSeason deletes a season from the db.
func (s *SeasonService) DeleteSeason(id mmm.SeasonID) error {
	// retrieve season.
	var v mmm.Season
	if err := s.client.db.One("id", id, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return ErrDatabaseQuery
	}

	// delete season.
	if err := s.client.db.DeleteStruct(&v); err != nil {
		s.client.logger.Log("err", ErrDatabaseDelete, "msg", err.Error())
		return ErrDatabaseDelete
	}
	return nil
}
