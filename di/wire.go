//go:build wireinject
// +build wireinject

// Package di Inject dependence by wire command.
package di

import (
	"github.com/google/wire"
	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/infrastructure/memory"
	"github.com/nao1215/sqluv/infrastructure/persistence"
	"github.com/nao1215/sqluv/interactor"
	"github.com/nao1215/sqluv/tui"
)

//go:generate wire

// New creates a new sqluv command instance.
func NewSqluv(arg *config.Argument) (*tui.TUI, func(), error) {
	wire.Build(
		config.Set,
		tui.Set,
		interactor.Set,
		persistence.Set,
		memory.Set,
	)
	return nil, nil, nil
}
