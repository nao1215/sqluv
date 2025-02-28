// Package interactor provides implementations of interfaces defined in the usecase.
package interactor

import "github.com/google/wire"

// Set is interactor providers.
var Set = wire.NewSet(
	NewFileReader,
)
