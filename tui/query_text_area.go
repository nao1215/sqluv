package tui

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// queryTextArea represents a query input field.
type queryTextArea struct {
	*tview.TextArea
	theme *Theme
}

// newQueryTextArea creates a new query input field.
func newQueryTextArea(theme *Theme) *queryTextArea {
	textArea := tview.NewTextArea().
		SetPlaceholder("Enter SQL query here...")

	textArea.SetBorder(true).
		SetTitle("Query").
		SetTitleAlign(tview.AlignLeft)

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

	q := &queryTextArea{
		TextArea: textArea,
		theme:    theme,
	}
	q.applyTheme(theme)
	return q
}

func (q *queryTextArea) applyTheme(theme *Theme) {
	q.theme = theme
	colors := theme.GetColors()

	// Update all text area colors
	q.SetBackgroundColor(colors.Background)
	q.SetTitleColor(colors.Foreground)
	q.SetTextStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))
	q.SetPlaceholderStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))

	// Update border color based on focus state
	if q.HasFocus() {
		q.SetBorderColor(colors.BorderFocus)
	} else {
		q.SetBorderColor(colors.Border)
	}
}
