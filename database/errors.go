package database

import "github.com/pmdcosta/mmm"

// Database errors.
const (
	ErrDatabaseFailed = mmm.Error("failed to start database client")
	ErrDatabaseInsert = mmm.Error("failed to insert data to the database")
	ErrDatabaseQuery  = mmm.Error("failed to query data from the database")
	ErrDatabaseDelete = mmm.Error("failed to delete data from the database")
	ErrDatabaseUpdate = mmm.Error("failed to update data from the database")
	ErrDatabaseMerge  = mmm.Error("failed to merge data from the database")
	ErrDatabaseExists = mmm.Error("already exists")
)
