package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// executeButton represents a button to execute SQL queries
type executeButton struct {
	*tview.Button
}

// newExecuteButton creates a new execute button
func newExecuteButton() *executeButton {
	button := tview.NewButton("Execute Query")
	button.SetBorder(false)
	button.SetTitle("Execute the SQL query")

	// Set text colors for better visibility
	button.SetLabelColor(tcell.ColorWhite)
	button.SetLabelColorActivated(tcell.ColorBlack)
	button.SetBackgroundColorActivated(tcell.ColorGreen)

	button.SetFocusFunc(func() {
		button.SetBorderColor(tcell.ColorGreen)
		button.SetLabelColor(tcell.ColorBlack)
		button.SetBackgroundColor(tcell.ColorLightGray)
	})

	button.SetBlurFunc(func() {
		button.SetBorderColor(tcell.ColorDefault)
		button.SetLabelColor(tcell.ColorWhite)
		button.SetBackgroundColor(tcell.ColorDefault)
	})

	return &executeButton{Button: button}
}
