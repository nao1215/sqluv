package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/rivo/tview"
)

// queryResultTable represents a table that displays SQL query results.
type queryResultTable struct {
	*tview.Table
	theme *Theme
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
	if table == nil {
		q.clear()
		return
	}
	q.Clear()

	// Get theme colors for consistent styling
	colors := q.theme.GetColors() // We need to store the theme in queryResultTable

	// Set column widths and headers
	for i, col := range table.Header() {
		q.SetCell(0, i,
			tview.NewTableCell(col).
				SetTextColor(colors.Header). // Use theme Header color
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
					SetTextColor(colors.Foreground). // Use theme Foreground color
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
