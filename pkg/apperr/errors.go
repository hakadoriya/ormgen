package apperr

import "errors"

var (
	ErrLanguageNotSupported     = errors.New("language not supported")
	ErrDialectNotSupported      = errors.New("dialect not supported")
	ErrInvalidArguments         = errors.New("invalid arguments")
	ErrSourcePathIsNotDirectory = errors.New("source path is not a directory")
	ErrNoSourceFound            = errors.New("no source found")
	ErrInvalidAnnotation        = errors.New("invalid annotation")
)
