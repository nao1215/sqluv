package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/usecase"
	"github.com/rivo/tview"
)

// TUI represents a text-based user interface.
type TUI struct {
	// files is a list of file paths that import to SQLite3 in-memory mode.
	files []*model.File

	// app is the TUI application.
	app *tview.Application
	// root is the root component of the TUI.
	root *home

	// usecases is a set of usecases.
	usecases *usecases
}

// home represents the home screen of the TUI.
type usecases struct {
	// fileReader is a usecase to read records from CSV/TSV/LTSV files and return them as model.Table.
	fileReader usecase.FileReader
}

// NewTUI creates a new TUI instance.
func NewTUI(
	arg *config.Argument,
	fileReader usecase.FileReader,
) *TUI {
	tui := &TUI{
		files: arg.Files(),
		root:  newHome(),
		app:   tview.NewApplication(),
		usecases: &usecases{
			fileReader: fileReader,
		},
	}
	tui.app.SetInputCapture(tui.keyBindings)
	tui.app.SetMouseCapture(tui.mouseHandler)
	tui.app.EnableMouse(true)
	tui.app.EnablePaste(true)

	return tui
}

// Run runs the TUI.
func (t *TUI) Run() error {
	return t.app.SetRoot(t.root.flex, true).Run()
}

// keyBindings handles key bindings.
func (t *TUI) keyBindings(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc, tcell.KeyCtrlC:
		t.app.Stop()
	}

	switch event.Rune() {
	case 'q':
		t.app.Stop()
	}
	return event
}

// mouseHandler handles mouse events.
// mouseHandler handles mouse events in the main text view.
func (t *TUI) mouseHandler(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	return event, action
}
