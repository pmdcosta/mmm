package mmm

// Season errors.
const (
	ErrSeasonRequired = Error("season required")
)

// Episode errors.
const (
	ErrEpisodeRequired = Error("episode required")
)

// Tag errors.
const (
	ErrTagRequired = Error("tag required")
)

// Media Type errors.
const (
	ErrTypeNotFound = Error("media type not found")
)

// Error represents an error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
