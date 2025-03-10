package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/rivo/tview"
)

// queryResultTable represents a table that displays SQL query results.
type queryResultTable struct {
	*tview.Table
}

// newQueryResultTable creates a new query result table.
func newQueryResultTable() *queryResultTable {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, true)

	table.SetTitle("Results").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(false)

	table.SetFixed(1, 0)

	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))
	queryResultTable := &queryResultTable{
		Table: table,
	}

	queryResultTable.SetFocusFunc(func() {
		queryResultTable.SetBorder(false)
		queryResultTable.SetBorderColor(tcell.ColorGreen)
		queryResultTable.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack))
	})

	queryResultTable.SetBlurFunc(func() {
		queryResultTable.SetBorderColor(tcell.ColorDefault)
		queryResultTable.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))
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
	if table == nil {
		q.clear()
		return
	}
	q.Clear()

	// Set column widths and headers
	for i, col := range table.Header() {
		q.SetCell(0, i,
			tview.NewTableCell(col).
				SetTextColor(tcell.ColorYellow).
				SetSelectable(false).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(q.calcMaxWidth(table)).
				SetExpansion(1))
	}

	// Set row data with consistent column widths
	rows := table.Records()
	for rowIdx, row := range rows {
		for colIdx, cell := range row {
			q.SetCell(rowIdx+1, colIdx,
				tview.NewTableCell(cell).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft).
					SetMaxWidth(q.calcMaxWidth(table)).
					SetSelectable(true).
					SetExpansion(1))
		}
	}
	q.ScrollToBeginning()

	// Set up row selection handling with the updated execution time
	q.setupRowSelectionHandling(stats, executionTime)
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
