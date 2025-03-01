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
	files    []*model.File      // list of file paths that import to SQLite3 in-memory mode.
	app      *tview.Application // TUI application.
	home     *home              // home component of the TUI.
	usecases *usecases          // set of usecases.
}

// home represents the home screen of the TUI.
type usecases struct {
	fileReader usecase.FileReader // read records from CSV/TSV/LTSV files and return them as model.Table.
}

// NewTUI creates a new TUI instance.
func NewTUI(
	arg *config.Argument,
	fileReader usecase.FileReader,
) *TUI {
	app := tview.NewApplication()
	tui := &TUI{
		files: arg.Files(),
		home:  newHome(app),
		app:   app,
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
	// Set the root of the application
	t.app.SetRoot(t.home.flex, true)

	t.home.footer.home()

	// Process files if available
	if t.hasLocalFiles() {
		// Load first file into the result table
		table, err := t.usecases.fileReader.Read(t.files[0])
		if err != nil {
			t.showError(err)
			return t.app.Run()
		}
		t.home.resultTable.update(table)
		// Add files to sidebar
		t.home.sidebar.addLocalFiles(t.files)

		// Setup sidebar selection handler
		t.home.sidebar.SetSelectedFunc(func(node *tview.TreeNode) {
			// If a file node is selected
			if file, ok := node.GetReference().(*model.File); ok {
				table, err := t.usecases.fileReader.Read(file)
				if err != nil {
					t.showError(err)
					return
				}
				t.home.resultTable.update(table)
			}
		})
	}
	return t.app.Run()
}

// hasLocalFiles returns true if there are local files.
func (t *TUI) hasLocalFiles() bool {
	return len(t.files) > 0
}

// showError displays an error dialog with the given message
func (t *TUI) showError(err error) {
	t.home.errorDialog.Show(t.home.flex, err.Error())
}

// keyBindings handles key bindings.
func (t *TUI) keyBindings(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc, tcell.KeyCtrlC:
		t.app.Stop()
	}
	return event
}

// mouseHandler handles mouse events.
func (t *TUI) mouseHandler(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	return event, action
}
