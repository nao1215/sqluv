// Package infrastructure manage sqluv infrastructure logic.
package infrastructure

import "errors"

var (
	// ErrNoLabel is error when label not found during LTSV parsing
	ErrNoLabel = errors.New("no labels in the data")
)
