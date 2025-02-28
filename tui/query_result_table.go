package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// queryResultTable represents a table that displays SQL query results.
type queryResultTable struct {
	table *tview.Table
}

// newQueryResultTable creates a new query result table.
func newQueryResultTable() *queryResultTable {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, true)

	table.SetTitle("Results").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	// Set header row style
	table.SetFixed(1, 0)
	table.SetSelectable(true, true)

	// Set default empty state
	displayEmptyState(table)

	return &queryResultTable{
		table: table,
	}
}

// displayEmptyState shows an empty state message in the table.
func displayEmptyState(table *tview.Table) {
	table.Clear()
	table.SetCell(0, 0, tview.NewTableCell("No query results to display").
		SetTextColor(tcell.ColorDefault).
		SetAlign(tview.AlignCenter).
		SetExpansion(1))
}

/*
// updateWithResults updates the table with query results.
// columns: Column names
// rows: Data rows (each row is a slice of string values)
func (q *queryResultTable) updateWithResults(columns []string, rows [][]string) {
	q.table.Clear()

	// Set column headers
	for i, col := range columns {
		q.table.SetCell(0, i,
			tview.NewTableCell(col).
				SetTextColor(tcell.ColorYellow).
				SetSelectable(false).
				SetAlign(tview.AlignCenter))
	}

	// Set data rows
	for rowIdx, row := range rows {
		for colIdx, cell := range row {
			q.table.SetCell(rowIdx+1, colIdx,
				tview.NewTableCell(cell).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft))
		}
	}

	// If no results, show empty state
	if len(rows) == 0 {
		displayEmptyState(q.table)
	}
}
*/
