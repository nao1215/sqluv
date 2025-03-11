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
func newHistoryButton(theme *Theme) *historyButton {
	button := tview.NewButton("History")
	button.SetBorder(false)
	button.SetTitle("Show SQL Query History")

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

	h := &historyButton{Button: button}
	h.applyTheme(theme)

	return h
}

func (h *historyButton) applyTheme(theme *Theme) {
	colors := theme.GetColors()

	h.SetBackgroundColor(colors.Button)
	h.SetLabelColor(colors.ButtonText)
	h.SetLabelColorActivated(colors.ButtonTextFocus)
	h.SetBackgroundColorActivated(colors.ButtonFocus)

	h.SetStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))

	if h.HasFocus() {
		h.SetBorderColor(colors.BorderFocus)
	} else {
		h.SetBorderColor(colors.Border)
	}
}
