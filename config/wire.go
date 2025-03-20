// Package config provides functionality to read environment variables,
// runtime arguments, and configuration files.
package config

import "github.com/google/wire"

// Set is config providers.
var Set = wire.NewSet(
	NewMemoryDB,
	NewDBConfig,
	NewColorConfig,
	NewHistoryDB,
)
