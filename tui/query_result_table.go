package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/rivo/tview"
)

// queryResultTable represents a table that displays SQL query results.
type queryResultTable struct {
	*tview.Table
	theme        *Theme
	columnOffset int // new field to track the starting column index
	maxColumns   int // new field to define how many columns to display
}

// newQueryResultTable creates a new query result table.
func newQueryResultTable(theme *Theme) *queryResultTable {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, true)

	table.SetTitle("Results").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(false)
	table.SetFixed(1, 0)

	// Use theme colors instead of hardcoded colors
	colors := theme.GetColors()
	table.SetTitleColor(colors.Foreground)
	table.SetBackgroundColor(colors.Background)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(colors.Selection).
		Foreground(colors.SelectionText))

	queryResultTable := &queryResultTable{
		Table: table,
		theme: theme,
	}

	queryResultTable.SetFocusFunc(func() {
		queryResultTable.SetBorder(false)
		colors := theme.GetColors() // Get fresh colors in case theme changed
		queryResultTable.SetBorderColor(colors.BorderFocus)
		queryResultTable.SetSelectedStyle(tcell.StyleDefault.
			Background(colors.Selection).
			Foreground(colors.SelectionText))
	})

	queryResultTable.SetBlurFunc(func() {
		colors := theme.GetColors() // Get fresh colors in case theme changed
		queryResultTable.SetBorderColor(colors.Border)
		queryResultTable.SetSelectedStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.Foreground))
	})

	queryResultTable.clear()
	return queryResultTable
}

// clear clears the table.
// After calling this method, the table will be empty and display an empty state message.
func (q *queryResultTable) clear() {
	q.Clear()
}

// update updates the table with model.Table data
func (q *queryResultTable) update(table *model.Table, stats *rowStatistics, executionTime float64) {
	q.Clear()
	colors := q.theme.GetColors()
	headers := table.Header()
	totalCols := len(headers)

	// Ensure that q.maxColumns is set to the desired number of visible columns.
	if q.maxColumns <= 0 {
		q.maxColumns = 6
	}

	// Adjust the column offset if needed.
	if q.columnOffset > totalCols-q.maxColumns {
		q.columnOffset = totalCols - q.maxColumns
		if q.columnOffset < 0 {
			q.columnOffset = 0
		}
	}

	// Calculate the ending index for visible columns.
	end := q.columnOffset + q.maxColumns
	if end > totalCols {
		end = totalCols
	}

	// Render header row for visible columns.
	for i, col := range headers[q.columnOffset:end] {
		q.SetCell(0, i,
			tview.NewTableCell(col).
				SetTextColor(colors.Header).
				SetSelectable(false).
				SetAlign(tview.AlignCenter))
	}

	// Render rows.
	rows := table.Records()
	for rowIdx, row := range rows {
		rEnd := q.columnOffset + q.maxColumns
		if rEnd > len(row) {
			rEnd = len(row)
		}
		for colIdx, cell := range row[q.columnOffset:rEnd] {
			q.SetCell(rowIdx+1, colIdx,
				tview.NewTableCell(cell).
					SetTextColor(colors.Foreground).
					SetAlign(tview.AlignLeft).
					SetMaxWidth(q.calcMaxWidth(table)).
					SetExpansion(1))
		}
	}

	q.ScrollToBeginning()
	q.setupRowSelectionHandling(stats, executionTime)
	q.setHorizontalScrollHandler(totalCols, table, stats, executionTime)
}

// setHorizontalScrollHandler sets the input capture to handle horizontal scrolling.
func (q *queryResultTable) setHorizontalScrollHandler(totalCols int, table *model.Table, stats *rowStatistics, executionTime float64) {
	q.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Get current selection (row and col in the visible area)
		_, selCol := q.GetSelection()
		switch event.Key() {
		case tcell.KeyLeft:
			// If focus is on the very first visible column and there are hidden columns to the left,
			// shift the columns only once.
			if selCol == 0 && q.columnOffset > 0 {
				q.columnOffset--
				q.update(table, stats, executionTime)
				return nil
			}
		case tcell.KeyRight:
			// If focus is on the very last visible column and there are hidden columns to the right,
			// shift the viewport.
			if selCol == q.maxColumns-1 && q.columnOffset+q.maxColumns < totalCols {
				q.columnOffset++
				q.update(table, stats, executionTime)
				return nil
			}
		}
		return event
	})
}

// calcMaxWidth calculates the maximum width of the table.
func (q *queryResultTable) calcMaxWidth(table *model.Table) int {
	columns := table.Header()
	columnCount := len(columns)

	_, _, w, _ := q.GetRect() //nolint:dogsled // x, y, width, height. only width is used
	colWidth := 0
	if columnCount > 0 {
		colWidth = (w - columnCount - 1) / columnCount // Account for borders
	}
	return colWidth
}

// Update the setupRowSelectionHandling function to track cell position
func (q *queryResultTable) setupRowSelectionHandling(stats *rowStatistics, executionTime float64) {
	q.SetSelectable(true, true) // Allow both row and column selection

	// Update selection handler to pass both row and column information
	q.SetSelectedFunc(func(row, col int) {
		// Don't count header row (row 0) as a data row
		if row > 0 {
			totalRows := q.GetRowCount() - 1 // Subtract 1 for header row
			stats.updateSelectedCell(row-1, col, totalRows, executionTime)
		}
	})

	// Add handler for when selection changes without clicking
	q.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 {
			totalRows := q.GetRowCount() - 1
			stats.updateSelectedCell(row-1, col, totalRows, executionTime)
		}
	})
}

func (q *queryResultTable) applyTheme(theme *Theme) {
	q.theme = theme
	colors := theme.GetColors()
	q.SetBackgroundColor(colors.Background)

	// Apply header row styling
	for i := 0; i < q.GetColumnCount(); i++ {
		cell := q.GetCell(0, i)
		if cell != nil {
			cell.SetTextColor(colors.Header)
			cell.SetBackgroundColor(colors.Background)
		}
	}

	// Update selection style
	if q.HasFocus() {
		q.SetBorderColor(colors.BorderFocus)
		q.SetSelectedStyle(tcell.StyleDefault.
			Background(colors.Selection).
			Foreground(colors.SelectionText))
	} else {
		q.SetBorderColor(colors.Border)
		q.SetSelectedStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.Foreground))
	}
}
