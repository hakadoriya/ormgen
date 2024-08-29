package apperr

import "errors"

var (
	ErrInvalidArguments         = errors.New("invalid arguments")
	ErrSourcePathIsNotDirectory = errors.New("source path is not a directory")
)
