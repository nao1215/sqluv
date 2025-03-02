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
	table.SetSelectable(true, true)

	queryResultTable := &queryResultTable{
		Table: table,
	}
	queryResultTable.clear()
	return queryResultTable
}

// clear clears the table.
// After calling this method, the table will be empty and display an empty state message.
func (q *queryResultTable) clear() {
	q.Clear()
	q.SetCell(0, 0, tview.NewTableCell("No query results to display").
		SetTextColor(tcell.ColorDefault).
		SetAlign(tview.AlignCenter).
		SetExpansion(1))
}

// update updates the table with model.Table data
func (q *queryResultTable) update(table *model.Table) {
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
					SetExpansion(1))
		}
	}
	q.ScrollToBeginning()
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
