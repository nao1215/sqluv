package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

	textArea.SetFocusFunc(func() {
		textArea.SetBorderColor(tcell.ColorGreen)
	})

	textArea.SetBlurFunc(func() {
		textArea.SetBorderColor(tcell.ColorDefault)
	})
	return &queryTextArea{
		TextArea: textArea,
	}
}
