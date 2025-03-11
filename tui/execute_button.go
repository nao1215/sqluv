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
func newExecuteButton(theme *Theme) *executeButton {
	button := tview.NewButton("Execute Query")
	button.SetBorder(false)
	button.SetTitle("Execute the SQL query")

	button.SetFocusFunc(func() {
		button.SetBorderColor(theme.GetColors().BorderFocus)
		button.SetLabelColor(theme.GetColors().ButtonTextFocus)
		button.SetBackgroundColor(theme.GetColors().ButtonFocus)
	})

	button.SetBlurFunc(func() {
		button.SetBorderColor(theme.GetColors().Border)
		button.SetLabelColor(theme.GetColors().ButtonText)
		button.SetBackgroundColor(theme.GetColors().Button)
	})

	b := &executeButton{Button: button}
	b.applyTheme(theme)

	return b
}

func (e *executeButton) applyTheme(theme *Theme) {
	colors := theme.GetColors()

	e.SetBackgroundColor(colors.Button)
	e.SetLabelColor(colors.ButtonText)
	e.SetLabelColorActivated(colors.ButtonTextFocus)
	e.SetBackgroundColorActivated(colors.ButtonFocus)
	e.SetStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))

	if e.HasFocus() {
		e.SetBorderColor(colors.BorderFocus)
	} else {
		e.SetBorderColor(colors.Border)
	}
}
