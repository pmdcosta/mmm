package database

import (
	"github.com/imdario/mergo"
	"github.com/pmdcosta/mmm"
)

// ensure SeriesService implements mmm.SeriesService.
var _ mmm.SeriesService = &SeriesService{}

// SeriesService series management service to interact with the database.
type SeriesService struct {
	client *Client
}

// ListSeries returns a list of series.
func (s *SeriesService) ListSeries() ([]mmm.Series, error) {
	// retrieve series.
	var v []mmm.Series
	if err := s.client.db.All(&v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return v, nil
}

// CreateSeries persists a series to the database.
func (s *SeriesService) CreateSeries(v *mmm.Series) (*mmm.Series, error) {
	// require object and id.
	if v == nil || v.Title == "" || v.Type == "" {
		return nil, mmm.ErrSeriesRequired
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

// Series retrieves a series from the database.
func (s *SeriesService) Series(id mmm.SeriesID) (*mmm.Series, error) {
	// retrieve series.
	var v mmm.Series
	if err := s.client.db.One("id", id, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return &v, nil
}

// UpdateSeries updates a series from the db.
func (s *SeriesService) UpdateSeries(id mmm.SeriesID, new *mmm.Series) (*mmm.Series, error) {
	// retrieve series.
	var old mmm.Series
	if err := s.client.db.One("id", id, &old); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}

	// merge series with old one.
	if err := mergo.Merge(&new, old); err != nil {
		s.client.logger.Log("err", ErrDatabaseMerge, "msg", err.Error())
		return nil, ErrDatabaseMerge
	}

	// update series.
	if err := s.client.db.Update(&new); err != nil {
		s.client.logger.Log("err", ErrDatabaseUpdate, "msg", err.Error())
		return nil, ErrDatabaseUpdate
	}
	return new, nil
}

// DeleteSeries deletes a series from the db.
func (s *SeriesService) DeleteSeries(id mmm.SeriesID) error {
	// retrieve series.
	var v mmm.Series
	if err := s.client.db.One("id", id, &v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return ErrDatabaseQuery
	}

	// delete series.
	if err := s.client.db.DeleteStruct(&v); err != nil {
		s.client.logger.Log("err", ErrDatabaseDelete, "msg", err.Error())
		return ErrDatabaseDelete
	}
	return nil
}
