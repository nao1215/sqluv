package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// historyButton represents a button to display SQL query history
type historyButton struct {
	*tview.Button
}

// newHistoryButton creates a new history button
func newHistoryButton() *historyButton {
	button := tview.NewButton("History")
	button.SetBorder(false)
	button.SetTitle("Show SQL Query History")

	// Set text colors for better visibility
	button.SetLabelColor(tcell.ColorWhite)
	button.SetLabelColorActivated(tcell.ColorBlack)
	button.SetBackgroundColorActivated(tcell.ColorGreen)

	button.SetFocusFunc(func() {
		button.SetBorderColor(tcell.ColorBlue)
		button.SetLabelColor(tcell.ColorBlack)
		button.SetBackgroundColor(tcell.ColorLightGray)
	})

	button.SetBlurFunc(func() {
		button.SetBorderColor(tcell.ColorDefault)
		button.SetLabelColor(tcell.ColorWhite)
		button.SetBackgroundColor(tcell.ColorDefault)
	})

	return &historyButton{Button: button}
}
