package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

// rowStatistics represents a component that displays row statistics
type rowStatistics struct {
	*tview.TextView
}

// newRowStatistics creates a new row statistics component
func newRowStatistics(theme *Theme) *rowStatistics {
	textView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	textView.SetBorder(false)
	// Set initial text
	textView.SetText("No rows to display")

	r := &rowStatistics{
		TextView: textView,
	}
	r.applyTheme(theme)
	return r
}

// clear resets the statistics display
func (r *rowStatistics) clear() {
	r.SetText("No rows to display")
}

// updateSelectedCell updates the display with the currently selected cell position
func (r *rowStatistics) updateSelectedCell(selectedRow, selectedCol int, totalRows int, executionTime float64) {
	var text string
	if totalRows == 0 {
		text = "No rows to display"
	} else if selectedRow == -1 {
		// No row selected yet
		text = fmt.Sprintf("[green]%d row(s) in set (%.3f sec)[white]", totalRows, executionTime)
	} else {
		// Show both position and total rows
		text = fmt.Sprintf("[green]Row %d, Column %d, total %d row(s) (%.3f sec)[white]",
			selectedRow+1, selectedCol+1, totalRows, executionTime)
	}
	r.SetText(text)
}

func (r *rowStatistics) applyTheme(theme *Theme) {
	colors := theme.GetColors()
	r.SetBackgroundColor(colors.Background)

	if r.HasFocus() {
		r.SetBorderColor(colors.BorderFocus)
	} else {
		r.SetBorderColor(colors.Border)
	}
}
