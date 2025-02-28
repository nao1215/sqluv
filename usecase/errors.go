package usecase

import "errors"

var (
	// ErrNotSupportedFileFormat is an error that occurs when the file format is not supported.
	ErrNotSupportedFileFormat = errors.New("not supported file format")
)
