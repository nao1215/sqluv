package tui

import "github.com/rivo/tview"

// queryTextArea represents a query input field.
type queryTextArea struct {
	*tview.TextArea
}

// newQueryTextArea creates a new query input field.
func newQueryTextArea() *queryTextArea {
	textArea := tview.NewTextArea().
		SetPlaceholder("Enter SQL query here...")
	textArea.SetBorder(true).
		SetTitle("Query").
		SetTitleAlign(tview.AlignLeft)

	return &queryTextArea{
		TextArea: textArea,
	}
}
