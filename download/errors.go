package download

import "github.com/pmdcosta/mmm"

const (
	ErrCreateFile   = mmm.Error("failed to create torrent file")
	ErrDownloadFile = mmm.Error("failed to download torrent file")
	ErrWriteFile    = mmm.Error("failed to write content to torrent file")
)
