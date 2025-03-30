// Package memory implements the memory-related interfaces defined in domain/repository.
package memory

import "github.com/google/wire"

// Set is memory wire set.
var Set = wire.NewSet(
	NewTableCreator,
	NewTableGetter,
	NewRecordInserter,
	NewQueryExecutor,
	NewStatementExecutor,
	NewTableDDLGetter,
)
