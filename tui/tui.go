package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/config"
	"github.com/rivo/tview"
)

// TUI represents a text-based user interface.
type TUI struct {
	// filePaths is a list of file paths that import to SQLite3 in-memory mode.
	filePaths []string

	// app is the TUI application.
	app *tview.Application
	// root is the root component of the TUI.
	root *home
}

// NewTUI creates a new TUI instance.
func NewTUI(
	arg *config.Argument,
) *TUI {
	tui := &TUI{
		filePaths: arg.FilePaths(),
		root:      newHome(),
		app:       tview.NewApplication(),
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
