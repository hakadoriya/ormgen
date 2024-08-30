package apperr

import "errors"

var (
	ErrLanguageNotSupported     = errors.New("language not supported")
	ErrInvalidArguments         = errors.New("invalid arguments")
	ErrSourcePathIsNotDirectory = errors.New("source path is not a directory")
	ErrNoStructSourceFound      = errors.New("no struct source found")
)
