package tui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
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
		fileWriter     usecase.FileWriter
		tableCreator   usecase.TableCreator
		tablesGetter   usecase.TablesGetter
		ddlGetter      usecase.TableDDLGetter
		sqlExecutor    usecase.SQLExecutor
		recordInserter usecase.RecordsInserter
	}

	// dbmsUsecases represents use cases for DBMS operations
	dbmsUsecases struct {
		queryExecutor usecase.QueryExecutor
		tablesGetter  usecase.TablesGetter
		ddlGetter     usecase.TableDDLInRemoteGetter
		fileWriter    usecase.FileWriter

		closeDB       func() // Added field for database cleanup function
		isDBConnected bool   // Flag to track if we're connected to a database
		databaseName  string // Name of the connected database
	}

	// historyUsecases represents use cases for history operations
	historyUsecases struct {
		historyTableCreator usecase.HistoryTableCreator
		historyCreator      usecase.HistoryCreator
		historyLister       usecase.HistoryLister
	}
)

// TUI represents a text-based user interface.
type TUI struct {
	files           []*model.File      // list of file paths that import to SQLite3 in-memory mode.
	app             *tview.Application // TUI application.
	home            *home              // home component of the TUI.
	localUsecases   *localUsecases
	dbmsUsecases    *dbmsUsecases
	historyUsecases *historyUsecases
	dbConfig        *config.DBConfig // Database configuration manager
	theme           *Theme

	lastExecutionTime float64      // Time taken to execute the last query
	latestTable       *model.Table // Latest table fetched from the database
}

// NewTUI creates a new TUI instance.
func NewTUI(
	arg *config.Argument,
	fileReader usecase.FileReader,
	fileWriter usecase.FileWriter,
	tableCreator usecase.TableCreator,
	tablesGetter usecase.TablesGetter,
	ddlGetter usecase.TableDDLGetter,
	sqlExecuter usecase.SQLExecutor,
	recordInserter usecase.RecordsInserter,
	historyTableCreator usecase.HistoryTableCreator,
	historyCreator usecase.HistoryCreator,
	historyLister usecase.HistoryLister,
	dbConfig *config.DBConfig,
	colorManager *config.ColorConfig,
) *TUI {
	app := tview.NewApplication()
	theme := NewTheme(colorManager, app)

	tui := &TUI{
		files: arg.Files(),
		home:  newHome(app, theme),
		app:   app,
		localUsecases: &localUsecases{
			fileReader:     fileReader,
			fileWriter:     fileWriter,
			tableCreator:   tableCreator,
			tablesGetter:   tablesGetter,
			ddlGetter:      ddlGetter,
			sqlExecutor:    sqlExecuter,
			recordInserter: recordInserter,
		},
		dbmsUsecases: &dbmsUsecases{
			fileWriter: fileWriter,
		},
		historyUsecases: &historyUsecases{
			historyTableCreator: historyTableCreator,
			historyCreator:      historyCreator,
			historyLister:       historyLister,
		},
		dbConfig: dbConfig,
		theme:    theme,
	}

	tui.app.SetInputCapture(tui.keyBindings)
	tui.app.SetMouseCapture(tui.mouseHandler)
	tui.app.EnableMouse(true)
	tui.app.EnablePaste(true)

	tui.home.applyTheme(theme)
	return tui
}

// Run runs the TUI.
func (t *TUI) Run() error {
	ctx := context.Background()
	t.app.SetRoot(t.home.flex, true)
	t.home.footer.setDefaulShortcut()

	if err := t.historyUsecases.historyTableCreator.CreateTable(ctx); err != nil {
		return fmt.Errorf("failed to create history table: %w", err)
	}

	if t.hasLocalFiles() {
		if err := t.importFiles(ctx); err != nil {
			return err
		}
	} else {
		connectionModal := newConnectionModal(t.app, t.theme, t.handleConnectionSelection)
		t.app.SetRoot(connectionModal.Modal, true)
		return t.app.Run()
	}

	t.app.SetFocus(t.home.queryTextArea)
	t.home.queryTextArea.applyTheme(t.theme)

	t.home.executeButton.SetSelectedFunc(func() {
		t.executeQuery(context.Background())
		t.app.SetFocus(t.home.queryTextArea)
	})
	t.home.historyButton.SetSelectedFunc(func() {
		t.showHistoryList()
	})
	t.refreshAllComponents()
	return t.app.Run()
}

// handleDBConnection is a generic function to handle database connections
func (t *TUI) handleDBConnection(conn *config.DBConnection) error {
	db, closeDB, err := t.newDatabaseConfig(conn)
	if err != nil {
		return err
	}

	// Initialize DBMS usecases
	queryExecutor := persistence.NewQueryExecutor(db)
	statementExecutor := persistence.NewStatementExecutor(db)
	tablesGetter := persistence.NewTablesGetter(db, conn)
	tableDDLGetter := persistence.NewTableDDLGetter(db, conn)

	t.dbmsUsecases = &dbmsUsecases{
		queryExecutor: interactor.NewQueryExecutor(queryExecutor, statementExecutor),
		tablesGetter:  interactor.NewTablesGetter(tablesGetter),
		ddlGetter:     interactor.NewTableDDLInRemoteGetter(tableDDLGetter),
	}

	// Store the database connection for later use
	t.dbmsUsecases.databaseName = conn.Database
	t.dbmsUsecases.closeDB = closeDB
	t.dbmsUsecases.isDBConnected = true

	// Load tables and update the sidebar
	t.loadDatabaseTables(context.Background(), conn.Database)

	// Successfully connected to the database
	t.app.SetRoot(t.home.flex, true)
	t.app.SetFocus(t.home.queryTextArea)
	t.home.queryTextArea.applyTheme(t.theme)

	return nil
}

// newDatabaseConfig creates a new database configuration
func (t *TUI) newDatabaseConfig(conn *config.DBConnection) (config.DBMS, func(), error) {
	var db config.DBMS
	var closeDB func()
	var err error

	switch conn.Type {
	case config.MySQL:
		mysqlConfig := config.NewMySQLConfig(
			conn.Host,
			conn.Port,
			conn.User,
			conn.Password,
			conn.Database,
		)
		db, closeDB, err = config.NewMySQLDB(mysqlConfig)
		if err != nil {
			return nil, nil, err
		}
	case config.PostgreSQL:
		pgConfig := config.NewPostgreSQLConfig(
			conn.Host,
			conn.Port,
			conn.User,
			conn.Password,
			conn.Database,
		)
		db, closeDB, err = config.NewPostgreSQLDB(pgConfig)
		if err != nil {
			return nil, nil, err
		}
	case config.SQLite3:
		sqliteConfig := config.NewSQLite3Config(conn.Database)
		db, closeDB, err = config.NewSQLite3DB(sqliteConfig)
		if err != nil {
			return nil, nil, err
		}
	case config.SQLServer:
		sqlserverConfig := config.NewSQLServerConfig(
			conn.Host,
			conn.Port,
			conn.User,
			conn.Password,
			conn.Database,
		)
		db, closeDB, err = config.NewSQLServerDB(sqlserverConfig)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("unsupported database type: %s", conn.Type)
	}
	return db, closeDB, nil
}

// handleConnectionSelection processes the selected database connection
func (t *TUI) handleConnectionSelection(conn *config.DBConnection) {
	if conn == nil {
		return // User canceled, exit the application
	}

	err := t.handleDBConnection(conn)
	if err != nil {
		// Connection failed, show error and offer to remove the connection
		t.showFailedConnectionDialog(conn, err)
		return
	}

	t.app.SetFocus(t.home.queryTextArea)
	t.home.queryTextArea.applyTheme(t.theme)

	t.home.executeButton.SetSelectedFunc(func() {
		t.executeQuery(context.Background())
		t.app.SetFocus(t.home.queryTextArea)
	})
	t.home.historyButton.SetSelectedFunc(func() {
		t.showHistoryList()
	})
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
			connectionModal := newConnectionModal(t.app, t.theme, t.handleConnectionSelection)
			t.app.SetRoot(connectionModal.Modal, true)
		})

	colors := t.theme.GetColors()
	errorModal.SetBorderStyle(tcell.StyleDefault.
		Foreground(colors.BorderFocus).
		Background(colors.Background))
	errorModal.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(colors.ButtonFocus).
		Foreground(colors.ButtonTextFocus))
	errorModal.SetButtonStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))
	errorModal.SetBackgroundColor(colors.Background)

	t.app.SetRoot(errorModal, true).SetFocus(errorModal)
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
	tables := []*model.Table{}
	defer func() {
		t.home.sidebar.update(tables, "local")
		t.home.queryTextArea.setTableNamesAndCloumnNamesCandidate(tables)
	}()

	for _, file := range t.files {
		table, err := t.localUsecases.fileReader.Read(ctx, file)
		if err != nil {
			return err
		}
		if err := t.localUsecases.tableCreator.CreateTable(ctx, table); err != nil {
			return err
		}
		if err := t.localUsecases.recordInserter.InsertRecords(ctx, table); err != nil {
			return err
		}
		tableList, err := t.localUsecases.tablesGetter.GetTables(ctx)
		if err != nil {
			return err
		}
		tables = tableList
	}
	return nil
}

// hasLocalFiles returns true if there are local files.
func (t *TUI) hasLocalFiles() bool {
	return len(t.files) > 0
}

// showError displays an error dialog with the given message
func (t *TUI) showError(err error) {
	t.home.dialog.Show(t.home.flex, "ERROR", err.Error())
}

func (t *TUI) keyBindings(event *tcell.EventKey) *tcell.EventKey {
	defer func() {
		if t.home.sidebar.HasFocus() && !t.home.footer.isActiveSearch() {
			t.home.footer.setSidebarShortcut()
			return
		}
		if !t.home.footer.isActiveSearch() {
			t.home.footer.setDefaulShortcut()
			return
		}
	}()

	// If sidebar has focus and "/" is pressed, activate footer search for sidebar fuzzy search.
	if t.home.sidebar.HasFocus() && event.Rune() == '/' {
		t.home.footer.ActivateSearch()
		t.home.footer.SetChangedFunc(func(text string) {
			t.home.sidebar.filterTables(text)
		})
		t.app.SetFocus(t.home.footer)
		return nil
	}

	switch {
	case event.Key() == tcell.KeyCtrlD:
		t.app.Stop()
	case event.Key() == tcell.KeyTAB:
		// Cycle focus: queryTextArea -> executeButton -> historyButton -> queryResultTable -> sidebar -> queryTextArea
		if t.home.queryTextArea.HasFocus() {
			t.app.SetFocus(t.home.executeButton)
			return nil
		} else if t.home.executeButton.HasFocus() {
			t.app.SetFocus(t.home.historyButton)
			return nil
		} else if t.home.historyButton.HasFocus() {
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
		// Reverse cycle: queryTextArea -> sidebar -> queryResultTable -> historyButton -> executeButton -> queryTextArea
		if t.home.queryTextArea.HasFocus() {
			t.app.SetFocus(t.home.sidebar)
			return nil
		} else if t.home.sidebar.HasFocus() {
			// If there are query results, focus on result table, otherwise skip to history button
			if t.hasQueryResults() {
				t.app.SetFocus(t.home.resultTable)
			} else {
				t.app.SetFocus(t.home.historyButton)
			}
			return nil
		} else if t.home.resultTable.HasFocus() {
			t.app.SetFocus(t.home.historyButton)
			return nil
		} else if t.home.historyButton.HasFocus() {
			t.app.SetFocus(t.home.executeButton)
			return nil
		} else if t.home.executeButton.HasFocus() {
			t.app.SetFocus(t.home.queryTextArea)
			return nil
		}
	case event.Key() == tcell.KeyCtrlT:
		t.theme.ShowColorSchemeSelector(func() {
			t.app.SetRoot(t.home.flex, true)
			t.app.SetFocus(t.home.queryTextArea)
			t.refreshAllComponents()
		})
		return nil
	case event.Key() == tcell.KeyCtrlS:
		t.showSaveDialog()
		return nil
	case event.Key() == tcell.KeyCtrlE:
		if t.home.queryTextArea.HasFocus() {
			t.executeQuery(context.Background())
			return nil
		}
		if t.home.sidebar.HasFocus() {
			node := t.home.sidebar.GetCurrentNode()
			if node != nil {
				if table, ok := node.GetReference().(*model.Table); ok {
					query := fmt.Sprintf("SELECT * FROM %s LIMIT 100", table.Name())
					t.home.queryTextArea.SetText(query, true)
					t.executeQuery(context.Background())
					return nil
				}
			}
		}
	case event.Key() == tcell.KeyCtrlH:
		t.showHistoryList()
		return nil
	case event.Key() == tcell.KeyF1:
		t.app.SetFocus(t.home.sidebar)
		return nil
	case event.Key() == tcell.KeyF2:
		t.app.SetFocus(t.home.queryTextArea)
		return nil
	case event.Key() == tcell.KeyF3:
		t.app.SetFocus(t.home.resultTable)
		return nil

	case event.Key() == tcell.KeyEscape:
		if t.home.footer.isActiveSearch() {
			t.home.footer.SetLabel("")
			t.home.sidebar.update(t.home.sidebar.allTables, t.home.sidebar.dbName)
			t.home.footer.DeactivateSearch()
			t.home.footer.update()
			t.app.SetFocus(t.home.sidebar)
		}
		return nil
	case event.Key() == tcell.KeyEnter:
		if !t.home.sidebar.HasFocus() {
			return event
		}
		node := t.home.sidebar.GetCurrentNode()
		if node == nil {
			return event
		}
		if table, ok := node.GetReference().(*model.Table); ok {
			ddlTables := []*model.Table{}
			if t.dbmsUsecases.isDBConnected {
				var err error
				ddlTables, err = t.dbmsUsecases.ddlGetter.GetTableDDL(context.Background(), table.Name())
				if err != nil {
					t.showError(err)
					return nil
				}
			} else {
				var err error
				ddlTables, err = t.localUsecases.ddlGetter.GetTableDDL(context.Background(), table.Name())
				if err != nil {
					t.showError(err)
					return nil
				}
			}
			if len(ddlTables) > 0 {
				// Display the first returned DDL table.
				t.home.resultTable.update(ddlTables[0], t.home.rowStatistics, 0)
			}
		}
		return nil
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

	if t.dbmsUsecases.isDBConnected && t.dbmsUsecases.queryExecutor != nil {
		err = t.executeDBMSQuery(ctx, sql)
	} else {
		err = t.executeLocalQuery(ctx, sql)
	}
	if err != nil {
		t.showError(fmt.Errorf("%w: sql='%s'", err, query))
		return
	}

	if err := t.recordUserRequest(ctx, query); err != nil {
		t.showError(fmt.Errorf("failed to record user request: %w", err))
		return
	}
}

// recordUserRequest record user request in DB.
func (t *TUI) recordUserRequest(ctx context.Context, request string) error {
	histories, err := t.historyUsecases.historyLister.List(ctx)
	if err != nil {
		return err
	}
	if err := t.historyUsecases.historyCreator.Create(ctx, model.NewHistory(len(histories)+1, request)); err != nil {
		return fmt.Errorf("failed to store user input history: %w", err)
	}
	return nil
}

// executeDBMSQuery executes SQL query against connected DBMS
func (t *TUI) executeDBMSQuery(ctx context.Context, sql *model.SQL) error {
	startTime := time.Now()
	output, err := t.dbmsUsecases.queryExecutor.ExecuteQuery(ctx, sql)
	if err != nil {
		return err
	}

	if sql.IsDDL() {
		t.loadDatabaseTables(ctx, t.dbmsUsecases.databaseName)
	}
	if sql.IsUpdate() {
		t.showRowsAffectedInfo(output.RowsAffected())
	}
	if output.HasTable() || sql.IsDelete() {
		t.lastExecutionTime = time.Since(startTime).Seconds()
		t.home.resultTable.update(output.Table(), t.home.rowStatistics, t.lastExecutionTime)
		t.updateRowStatistics(output.Table(), startTime)
		t.latestTable = output.Table()
	}
	return nil
}

// executeLocalQuery executes SQL query against local file data
func (t *TUI) executeLocalQuery(ctx context.Context, sql *model.SQL) error {
	startTime := time.Now()
	output, err := t.localUsecases.sqlExecutor.ExecuteSQL(ctx, sql)
	if err != nil {
		return err
	}

	if sql.IsUpdate() {
		t.showRowsAffectedInfo(output.RowsAffected())
	}
	if output.HasTable() || sql.IsDelete() {
		t.lastExecutionTime = time.Since(startTime).Seconds()
		t.home.resultTable.update(output.Table(), t.home.rowStatistics, t.lastExecutionTime)
		t.updateRowStatistics(output.Table(), startTime)
		t.latestTable = output.Table()
	}
	return nil
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
	if !t.dbmsUsecases.isDBConnected || t.dbmsUsecases.tablesGetter == nil {
		return
	}

	// Get tables from the database
	tables, err := t.dbmsUsecases.tablesGetter.GetTables(ctx)
	if err != nil {
		t.showError(fmt.Errorf("failed to load tables: %w", err))
		return
	}
	t.home.sidebar.update(tables, dbName)
	t.home.queryTextArea.setTableNamesAndCloumnNamesCandidate(tables)
}

// showHistoryList displays a list of SQL query history and allows selection
func (t *TUI) showHistoryList() {
	histories, err := t.historyUsecases.historyLister.List(context.Background())
	if err != nil {
		t.showError(err)
		return
	}
	if len(histories) == 0 {
		t.showError(errors.New("no query history available"))
		return
	}

	colors := t.theme.GetColors()
	// Create a list to display history items
	list := tview.NewList()
	list.SetTitle("SQL Query History")
	list.SetBorder(true)
	list.SetTitleColor(colors.Header)
	list.SetBorderStyle(tcell.StyleDefault.
		Background(colors.Border).
		Foreground(colors.BorderFocus))
	list.SetBackgroundColor(colors.Background)
	list.SetBorderStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.BorderFocus))
	list.SetMainTextStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))
	list.SetSelectedStyle(tcell.StyleDefault.
		Background(colors.BorderFocus).
		Foreground(colors.Foreground))

	// updateList refreshes the list using a fuzzy search query.
	updateList := func(query string) {
		list.Clear()
		// Add history items in reverse order (newest first)
		for i := len(histories) - 1; i >= 0; i-- {
			history := histories[i]
			displayText := normalizeSpaces(history.Request)

			// If query is non-empty, filter with fuzzy matching
			if query != "" && !fuzzy.Match(query, displayText) {
				continue
			}

			if len(displayText) > 75 {
				displayText = displayText[:71] + "..."
			}
			// Use closure to capture the correct history item.
			func(h model.History) {
				list.AddItem(displayText, "", 0, func() {
					t.home.queryTextArea.SetText(h.Request, true)
					t.app.SetRoot(t.home.flex, true)
					t.app.SetFocus(t.home.queryTextArea)
				})
			}(history)
		}
	}

	// Initial population of the list without a filter.
	updateList("")
	// Create an input field for fuzzy search.
	searchInput := tview.NewInputField().
		SetLabel("Fuzzy Search: ").
		SetFieldStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.Foreground)).
		SetLabelStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.Header)).
		SetFieldTextColor(colors.SelectionText).
		SetFieldWidth(0).
		SetChangedFunc(func(text string) {
			updateList(text)
		})
	searchInput.SetTitleColor(colors.Header)
	searchInput.SetBackgroundColor(colors.Background)
	searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			t.app.SetFocus(list)
			return nil
		}
		return event
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			t.app.SetFocus(searchInput)
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			t.app.SetRoot(t.home.flex, true)
			t.app.SetFocus(t.home.queryTextArea)
			t.refreshAllComponents()
			return nil
		}
		return event
	})

	// Create a layout using a vertical flex: search input on top, list below.
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchInput, 1, 0, true).
		AddItem(list, 0, 1, false)

	// Optional: add a cancel behavior when user finishes editing.
	list.SetDoneFunc(func() {
		t.app.SetRoot(t.home.flex, true)
		t.app.SetFocus(t.home.queryTextArea)
	})
	t.app.SetRoot(flex, true)
}

var spaceRegex = regexp.MustCompile(`\s+`)

func normalizeSpaces(input string) string {
	input = strings.ReplaceAll(input, "\r\n", " ")
	input = strings.ReplaceAll(input, "\n", " ")
	input = spaceRegex.ReplaceAllString(input, " ")
	return strings.TrimSpace(input)
}

// updateRowStatistics updates the row statistics component with the result information
func (t *TUI) updateRowStatistics(table *model.Table, startTime time.Time) {
	if table == nil {
		t.home.rowStatistics.clear()
		return
	}

	t.lastExecutionTime = time.Since(startTime).Seconds()
	rowCount := len(table.Records())
	// Use -1 for row and column to indicate no selection yet
	t.home.rowStatistics.updateSelectedCell(-1, -1, rowCount, t.lastExecutionTime)
}

// refreshAllComponents new method to refresh all components with the current theme
func (t *TUI) refreshAllComponents() {
	t.home.applyTheme(t.theme)
}

// showSaveDialog displays an input form for the file path.
func (t *TUI) showSaveDialog() {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	if t.latestTable == nil {
		t.showError(errors.New("no query results to save"))
		return
	}

	colors := t.theme.GetColors()

	form := tview.NewForm()
	form.AddInputField("Save File Path", cwd, 0, nil, nil).
		AddButton("Save", func() {
			filePath := form.GetFormItem(0).(*tview.InputField).GetText()
			f, err := model.NewFile(filePath)
			if err != nil {
				t.showError(fmt.Errorf("failed to create file handle: %w", err))
				return
			}
			if err := t.localUsecases.fileWriter.WriteFile(context.Background(), f, t.latestTable); err != nil {
				t.showError(fmt.Errorf("failed to write file: %w", err))
				return
			}
			t.home.dialog.Show(t.home.flex, " 🚀 ", "File saved successfully")
		}).
		AddButton("Cancel", func() {
			t.app.SetRoot(t.home.flex, true)
		})

	form.SetBorder(true).
		SetTitle("Save Query Results").
		SetTitleAlign(tview.AlignCenter).
		SetBorderStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.BorderFocus))
	form.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(colors.ButtonFocus).
		Foreground(colors.ButtonTextFocus)).
		SetButtonDisabledStyle(tcell.StyleDefault.
			Background(colors.Button).
			Foreground(colors.ButtonText)).
		SetButtonStyle(tcell.StyleDefault.
			Background(colors.Button).
			Foreground(colors.ButtonText)).
		SetFieldStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.Foreground)).
		SetBackgroundColor(colors.Background)
	t.app.SetRoot(form, true)
}
