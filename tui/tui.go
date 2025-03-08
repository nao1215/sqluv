package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/infrastructure/persistence"
	"github.com/nao1215/sqluv/interactor"
	"github.com/nao1215/sqluv/usecase"
	"github.com/rivo/tview"
)

// Split usecases into two separate structures
type (
	// localUsecases represents use cases for local file operations
	localUsecases struct {
		fileReader     usecase.FileReader
		tableCreator   usecase.TableCreator
		sqlExecutor    usecase.SQLExecutor
		recordInserter usecase.RecordsInserter
	}

	// dbmsUsecases represents use cases for DBMS operations
	dbmsUsecases struct {
		queryExecutor usecase.QueryExecutor
		tablesGetter  usecase.TablesGetter
	}
)

// TUI represents a text-based user interface.
type TUI struct {
	files         []*model.File      // list of file paths that import to SQLite3 in-memory mode.
	app           *tview.Application // TUI application.
	home          *home              // home component of the TUI.
	localUsecases *localUsecases
	dbmsUsecases  *dbmsUsecases

	closeDB       func() // Added field for database cleanup function
	isDBConnected bool   // Flag to track if we're connected to a database
	databaseName  string // Name of the connected database
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
		localUsecases: &localUsecases{
			fileReader:     fileReader,
			tableCreator:   tableCreator,
			sqlExecutor:    sqlExecuter,
			recordInserter: recordInserter,
		},
		dbmsUsecases: &dbmsUsecases{},
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
	t.home.footer.setDefaulShortcut()

	if t.hasLocalFiles() {
		if err := t.importFiles(ctx); err != nil {
			t.showError(err)
			t.app.SetFocus(t.home.queryTextArea)
			return t.app.Run()
		}
	} else {
		connectionModal := newConnectionModal(t.app, t.handleConnectionSelection)
		t.app.SetRoot(connectionModal.Modal, true)
		return t.app.Run()
	}

	t.app.SetFocus(t.home.queryTextArea)
	t.home.executeButton.SetSelectedFunc(func() {
		t.executeQuery(context.Background())
	})
	return t.app.Run()
}

// handleConnectionSelection processes the selected database connection
func (t *TUI) handleConnectionSelection(conn *config.DBConnection) {
	if conn == nil {
		return // User canceled, exit the application
	}

	// Connect to the database based on connection type
	switch conn.Type {
	case config.MySQL:
		t.handleMySQLConnection(conn)
	case config.PostgreSQL:
		t.handlePostgreSQLConnection(conn)
	default:
		t.showError(fmt.Errorf("unsupported database type: %s", conn.Type))
	}

	t.app.SetFocus(t.home.queryTextArea)
	t.home.executeButton.SetSelectedFunc(func() {
		t.executeQuery(context.Background())
	})
}

// Add new method for handling PostgreSQL connections
func (t *TUI) handlePostgreSQLConnection(conn *config.DBConnection) {
	pgConfig := config.NewPostgreSQLConfig(
		conn.Host,
		conn.Port,
		conn.User,
		conn.Password,
		conn.Database,
	)

	db, closeDB, err := config.NewPostgreSQLDB(pgConfig)
	if err != nil {
		// Connection failed, show error and offer to remove the connection
		t.showFailedConnectionDialog(conn, err)
		return
	}

	// Store the database connection for later use
	t.databaseName = conn.Database
	t.closeDB = closeDB
	t.isDBConnected = true

	// Initialize DBMS usecases
	queryExecutor := persistence.NewQueryExecutor(db)
	statementExecutor := persistence.NewStatementExecutor(db)
	tablesGetter := persistence.NewTablesGetter(db, conn.Database, conn.Type)

	t.dbmsUsecases = &dbmsUsecases{
		queryExecutor: interactor.NewQueryExecutor(queryExecutor, statementExecutor),
		tablesGetter:  interactor.NewTablesGetter(tablesGetter),
	}

	// Load tables and update the sidebar
	t.loadDatabaseTables(context.Background(), conn.Database)

	// Successfully connected to the database
	t.app.SetRoot(t.home.flex, true)
	t.app.SetFocus(t.home.queryTextArea)
}

// Update the MySQL connection handler for consistency
func (t *TUI) handleMySQLConnection(conn *config.DBConnection) {
	mysqlConfig := config.NewMySQLConfig(
		conn.Host,
		conn.Port,
		conn.User,
		conn.Password,
		conn.Database,
	)

	db, closeDB, err := config.NewMySQLDB(mysqlConfig)
	if err != nil {
		// Connection failed, show error and offer to remove the connection
		t.showFailedConnectionDialog(conn, err)
		return
	}

	// Store the database connection for later use
	t.databaseName = conn.Database
	t.closeDB = closeDB
	t.isDBConnected = true

	// Initialize DBMS usecases
	queryExecutor := persistence.NewQueryExecutor(db)
	statementExecutor := persistence.NewStatementExecutor(db)
	tablesGetter := persistence.NewTablesGetter(db, conn.Database, conn.Type)

	t.dbmsUsecases = &dbmsUsecases{
		queryExecutor: interactor.NewQueryExecutor(queryExecutor, statementExecutor),
		tablesGetter:  interactor.NewTablesGetter(tablesGetter),
	}

	// Load tables and update the sidebar
	t.loadDatabaseTables(context.Background(), conn.Database)

	// Successfully connected to the database
	t.app.SetRoot(t.home.flex, true)
	t.app.SetFocus(t.home.queryTextArea)
}

// showFailedConnectionDialog shows a dialog for failed connections with option to remove
func (t *TUI) showFailedConnectionDialog(conn *config.DBConnection, err error) {
	errorModal := tview.NewModal().
		SetText(fmt.Sprintf("Failed to connect to database: %v\nDo you want to remove this connection from saved configurations?", err)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "Yes" {
				t.removeConnectionFromConfig(conn.Name)
			}
			// Return to the connection modal
			connectionModal := newConnectionModal(t.app, t.handleConnectionSelection)
			t.app.SetRoot(connectionModal.Modal, true)
		})
	t.app.SetRoot(errorModal, true)
}

// removeConnectionFromConfig removes a connection from the saved configurations
func (t *TUI) removeConnectionFromConfig(connectionName string) {
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		t.showError(fmt.Errorf("could not initialize config manager: %w", err))
		return
	}

	err = dbConfig.RemoveConnection(connectionName)
	if err != nil {
		t.showError(fmt.Errorf("failed to remove connection: %w", err))
	}
}

// importFiles imports files into the SQLite3 in-memory database.
func (t *TUI) importFiles(ctx context.Context) error {
	importedFiles := make([]*model.File, 0, len(t.files))
	defer func() {
		t.home.sidebar.updateLocalFiles(importedFiles)
	}()

	for _, file := range t.files {
		table, err := t.localUsecases.fileReader.Read(file)
		if err != nil {
			return err
		}
		if err := t.localUsecases.tableCreator.CreateTable(ctx, table); err != nil {
			return err
		}
		if err := t.localUsecases.recordInserter.InsertRecords(ctx, table); err != nil {
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

func (t *TUI) keyBindings(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Key() == tcell.KeyEsc, event.Key() == tcell.KeyCtrlD:
		t.app.Stop()
	case event.Key() == tcell.KeyTAB:
		// Cycle focus: queryTextArea -> executeButton -> queryResultTable -> sidebar -> queryTextArea
		if t.home.queryTextArea.HasFocus() {
			t.app.SetFocus(t.home.executeButton)
			return nil
		} else if t.home.executeButton.HasFocus() {
			// If there are query results, focus on result table, otherwise skip to sidebar
			if t.hasQueryResults() {
				t.app.SetFocus(t.home.resultTable)
			} else {
				t.app.SetFocus(t.home.sidebar)
			}
			return nil
		} else if t.home.resultTable.HasFocus() {
			t.app.SetFocus(t.home.sidebar)
			return nil
		} else if t.home.sidebar.HasFocus() {
			t.app.SetFocus(t.home.queryTextArea)
			return nil
		}
	case event.Key() == tcell.KeyBacktab: // SHIFT+TAB
		// Reverse cycle: queryTextArea -> sidebar -> queryResultTable -> executeButton -> queryTextArea
		if t.home.queryTextArea.HasFocus() {
			t.app.SetFocus(t.home.sidebar)
			return nil
		} else if t.home.sidebar.HasFocus() {
			// If there are query results, focus on result table, otherwise skip to execute button
			if t.hasQueryResults() {
				t.app.SetFocus(t.home.resultTable)
			} else {
				t.app.SetFocus(t.home.executeButton)
			}
			return nil
		} else if t.home.resultTable.HasFocus() {
			t.app.SetFocus(t.home.executeButton)
			return nil
		} else if t.home.executeButton.HasFocus() {
			t.app.SetFocus(t.home.queryTextArea)
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
	if t.isDBConnected && t.dbmsUsecases.queryExecutor != nil {
		t.executeDBMSQuery(ctx, sql, query)
	} else {
		t.executeLocalQuery(ctx, sql, query)
	}
}

// executeDBMSQuery executes SQL query against connected DBMS
func (t *TUI) executeDBMSQuery(ctx context.Context, sql *model.SQL, rawQuery string) {
	output, err := t.dbmsUsecases.queryExecutor.ExecuteQuery(ctx, sql)
	if err != nil {
		t.showError(fmt.Errorf("%w: sql='%s'", err, rawQuery))
		return
	}

	if sql.IsDDL() {
		t.loadDatabaseTables(ctx, t.databaseName)
	}
	if sql.IsUpdate() {
		t.showRowsAffectedInfo(output.RowsAffected())
	}
	if output.HasTable() || sql.IsDelete() {
		t.home.resultTable.update(output.Table())
	}
}

// executeLocalQuery executes SQL query against local file data
func (t *TUI) executeLocalQuery(ctx context.Context, sql *model.SQL, rawQuery string) {
	output, err := t.localUsecases.sqlExecutor.ExecuteSQL(ctx, sql)
	if err != nil {
		t.showError(fmt.Errorf("%w: sql='%s'", err, rawQuery))
		return
	}

	if sql.IsUpdate() {
		t.showRowsAffectedInfo(output.RowsAffected())
	}
	if output.HasTable() || sql.IsDelete() {
		t.home.resultTable.update(output.Table())
	}
}

// showRowsAffectedInfo displays information about rows affected by a DML operation
func (t *TUI) showRowsAffectedInfo(rowsAffected int64) {
	infoModal := tview.NewModal().
		SetText(fmt.Sprintf("%d row(s) affected", rowsAffected)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			// Return to the main UI
			t.app.SetRoot(t.home.flex, true)
			t.app.SetFocus(t.home.queryTextArea)
		})

	// Show the modal over the current UI
	pages := tview.NewPages().
		AddPage("background", t.home.flex, true, true).
		AddPage("modal", infoModal, true, true)

	t.app.SetRoot(pages, true)
	t.app.SetFocus(infoModal)
}

// hasQueryResults checks if the result table has any query results
func (t *TUI) hasQueryResults() bool {
	// Check if there's more than just the header row (the "No query results" message)
	rowCount := t.home.resultTable.GetRowCount()
	if rowCount <= 1 {
		return false
	}

	// Alternative check: Check if the first cell doesn't contain the "No query results" message
	cell := t.home.resultTable.GetCell(0, 0)
	return cell != nil && cell.Text != "No query results to display"
}

// mouseHandler handles mouse events.
func (t *TUI) mouseHandler(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	return event, action
}

// loadDatabaseTables fetches tables from the database and updates the sidebar
func (t *TUI) loadDatabaseTables(ctx context.Context, dbName string) {
	if !t.isDBConnected || t.dbmsUsecases.tablesGetter == nil {
		return
	}

	// Get tables from the database
	tables, err := t.dbmsUsecases.tablesGetter.GetTables(ctx)
	if err != nil {
		t.showError(fmt.Errorf("failed to load tables: %w", err))
		return
	}

	// Update the sidebar with the tables
	t.home.sidebar.updateTables(tables, dbName)
}
