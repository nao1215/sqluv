// Package tui manage text user interface.
package tui

import "github.com/google/wire"

// Set is shell wire set.
var Set = wire.NewSet(
	NewTUI,
)
