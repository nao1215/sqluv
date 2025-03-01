//go:build wireinject
// +build wireinject

// Package di Inject dependence by wire command.
package di

import (
	"github.com/google/wire"
	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/infrastructure/persistence"
	"github.com/nao1215/sqluv/interactor"
	"github.com/nao1215/sqluv/tui"
)

//go:generate wire

// New creates a new sqluv command instance.
func NewSqluv(arg *config.Argument) (*tui.TUI, error) {
	wire.Build(
		tui.Set,
		interactor.Set,
		persistence.Set,
	)
	return nil, nil
}
