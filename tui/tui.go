package tui

import (
	"context"
	"fmt"

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
	fileReader     usecase.FileReader      // read records from CSV/TSV/LTSV files and return them as model.Table.
	tableCreator   usecase.TableCreator    // create a table in memory database.
	sqlExecutor    usecase.SQLExecutor     // execute a query in memory database.
	recordInserter usecase.RecordsInserter // insert records in memory database.
}

// NewTUI creates a new TUI instance.
func NewTUI(
	arg *config.Argument,
	fileReader usecase.FileReader,
	tableCreator usecase.TableCreator,
	sqlExecuter usecase.SQLExecutor,
	recordInserter usecase.RecordsInserter,
) *TUI {
	app := tview.NewApplication()
	tui := &TUI{
		files: arg.Files(),
		home:  newHome(app),
		app:   app,
		usecases: &usecases{
			fileReader:     fileReader,
			tableCreator:   tableCreator,
			sqlExecutor:    sqlExecuter,
			recordInserter: recordInserter,
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
	ctx := context.Background()
	t.app.SetRoot(t.home.flex, true)
	t.home.footer.home()

	if t.hasLocalFiles() {
		if err := t.importFiles(ctx); err != nil {
			t.showError(err)
			return t.app.Run()
		}
	}
	return t.app.Run()
}

// importFiles imports files into the SQLite3 in-memory database.
func (t *TUI) importFiles(ctx context.Context) error {
	importedFiles := make([]*model.File, 0, len(t.files))
	defer func() {
		t.home.sidebar.addLocalFiles(importedFiles)
	}()

	for _, file := range t.files {
		table, err := t.usecases.fileReader.Read(file)
		if err != nil {
			return err
		}
		if err := t.usecases.tableCreator.CreateTable(ctx, table); err != nil {
			return err
		}
		if err := t.usecases.recordInserter.InsertRecords(ctx, table); err != nil {
			return err
		}
		importedFiles = append(importedFiles, file)
	}
	return nil
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
	case tcell.KeyTab:
		if t.home.queryTextArea.HasFocus() {
			t.executeQuery(context.Background())
			return nil
		}
	}
	return event
}

// executeQuery executes the SQL query in the query text area
func (t *TUI) executeQuery(ctx context.Context) {
	query := t.home.queryTextArea.GetText()
	sql, err := model.NewSQL(query)
	if err != nil {
		t.showError(err)
		return
	}

	output, err := t.usecases.sqlExecutor.ExecuteSQL(ctx, sql)
	if err != nil {
		t.showError(fmt.Errorf("%w: sql='%s'", err, query))
		return
	}

	if output.HasTable() {
		t.home.resultTable.update(output.Table())
	}
	// TODO: Show rows affected in the footer
}

// mouseHandler handles mouse events.
func (t *TUI) mouseHandler(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	return event, action
}
