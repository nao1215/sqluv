package tui

import (
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/rivo/tview"
)

// queryTextArea represents a query input field.
type queryTextArea struct {
	*tview.TextArea
	theme          *Theme
	candidates     []string // SQL reserved words, table names, and column names.
	completionList *tview.List
}

// newQueryTextArea creates a new query input field.
func newQueryTextArea(theme *Theme) *queryTextArea {
	textArea := tview.NewTextArea().
		SetPlaceholder("Enter SQL query here...")

	textArea.SetBorder(true).
		SetTitle("Query").
		SetTitleAlign(tview.AlignLeft)

	q := &queryTextArea{
		TextArea: textArea,
		theme:    theme,
	}
	q.candidates = q.sqlCandidate()

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

		// Auto-completion: if the user presses the right arrow key, we will try to auto-complete the current word.
		if event.Key() == tcell.KeyRight {
			text := textArea.GetText()
			words := strings.Fields(strings.ReplaceAll(text, "\n", " "))
			currentWord := words[len(words)-1]

			// Look for candidate completions.
			var match string
			for _, candidate := range q.candidates {
				if strings.HasPrefix(strings.ToUpper(candidate), strings.ToUpper(currentWord)) {
					match = candidate
					break
				}
			}
			// If a match is found, replace the current word with the full candidate.
			if match != "" && match != currentWord {
				start := len(text)
				for i := start - 1; i >= 0; i-- {
					if text[i] == ' ' || text[i] == '\n' || text[i] == '\t' {
						break
					}
					start = i
				}
				textArea.SetText(text[:start]+match, true)
				return nil
			}
		}
		return event // Pass other keys through
	})

	q.applyTheme(theme)
	return q
}

func (q *queryTextArea) setTableNamesAndCloumnNamesCandidate(tables []*model.Table) {
	q.candidates = q.sqlCandidate() // clear previous candidates
	for _, table := range tables {
		q.candidates = append(q.candidates, table.Name())
		for _, col := range table.Header() {
			q.candidates = append(q.candidates, col)
		}
	}
}

func (q *queryTextArea) sqlCandidate() []string {
	return []string{"SELECT", "FROM", "WHERE", "INSERT", "UPDATE", "DELETE", "JOIN", "GROUP", "ORDER", "LIMIT", "WITH"}
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
