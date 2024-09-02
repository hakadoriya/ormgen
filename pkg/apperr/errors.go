package apperr

import "errors"

var (
	ErrLanguageNotSupported     = errors.New("language not supported")
	ErrEmpty                    = errors.New("empty")
	ErrInvalidArguments         = errors.New("invalid arguments")
	ErrSourcePathIsNotDirectory = errors.New("source path is not a directory")
	ErrNoSourceFound            = errors.New("no source found")
)
