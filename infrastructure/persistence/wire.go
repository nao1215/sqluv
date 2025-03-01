// Package persistence implements the persistence-related interfaces defined in domain/repository.
package persistence

import (
	"github.com/google/wire"
)

// Set is persistence providers.
var Set = wire.NewSet(
	NewCSVReader,
	NewTSVReader,
	NewLTSVReader,
)
