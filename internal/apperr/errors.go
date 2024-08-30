package apperr

import "errors"

var (
	ErrInvalidArguments                = errors.New("invalid arguments")
	ErrSourcePathIsNotDirectory        = errors.New("source path is not a directory")
	ErrFailedToDetectPackageImportPath = errors.New("failed to detect package import path; Please use option, or run include the package in your GOPATH or module (GO111MODULE=auto may be required)")
	ErrNoStructSourceFound             = errors.New("no struct source found")
)
