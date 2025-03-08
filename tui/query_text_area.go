package tui

import (
	"github.com/atotto/clipboard"
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

	// Add keyboard shortcut handling for copy/paste
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Check for Ctrl+C (copy)
		if event.Key() == tcell.KeyCtrlC {
			_, from, to := textArea.GetSelection()
			if from != to {
				selectedText := textArea.GetText()
				if err := clipboard.WriteAll(selectedText[from:to]); err != nil {
					return event // Pass through
				}
				return nil // Consume the event
			}
			return event // No selection, pass through
		}

		// Check for Ctrl+V (paste)
		if event.Key() == tcell.KeyCtrlV {
			text, err := clipboard.ReadAll()
			if err != nil {
				return event // Pass through
			}
			textArea.SetText(text, true)
			return nil // Consume the event
		}

		// Check for Ctrl+X (cut)
		if event.Key() == tcell.KeyCtrlX {
			_, from, to := textArea.GetSelection()
			if from != to {
				selectedText := textArea.GetText()
				if err := clipboard.WriteAll(selectedText[from:to]); err != nil {
					return event // Pass through
				}
				textArea.Replace(from, to, "")
				return nil // Consume the event
			}
			return event // No selection, pass through
		}
		return event // Pass other keys through
	})

	return &queryTextArea{
		TextArea: textArea,
	}
}
