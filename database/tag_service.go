package database

import (
	"github.com/pmdcosta/mmm"
)

// ensure TagService implements mmm.TagService.
var _ mmm.TagService = &TagService{}

// TagService tag management service to interact with the database.
type TagService struct {
	client *Client
}

// ListTags returns a list of tags.
func (s *TagService) ListTags() ([]mmm.Tag, error) {
	// retrieve tags.
	var v []mmm.Tag
	if err := s.client.db.All(&v); err != nil {
		s.client.logger.Log("err", ErrDatabaseQuery, "msg", err.Error())
		return nil, ErrDatabaseQuery
	}
	return v, nil
}

// CreateTag persists a tag to the database.
func (s *TagService) CreateTag(v *mmm.Tag) error {
	// require id.
	if v.ID == "" {
		return mmm.ErrTagRequired
	}

	// save record.
	if err := s.client.db.Save(v); err != nil {
		s.client.logger.Log("err", ErrDatabaseInsert, "msg", err.Error())
		return ErrDatabaseInsert
	}
	return nil
}

// DeleteTag deletes a tag from the db.
func (s *TagService) DeleteTag(v *mmm.Tag) error {
	// delete tag.
	if err := s.client.db.DeleteStruct(&v); err != nil {
		s.client.logger.Log("err", ErrDatabaseDelete, "msg", err.Error())
		return ErrDatabaseDelete
	}
	return nil
}
