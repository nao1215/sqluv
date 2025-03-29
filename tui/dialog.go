package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// dialog represents a modal dialog for displaying error messages, or other.
type dialog struct {
	*tview.Modal
	previousPage tview.Primitive
	app          *tview.Application
}

// newDialog creates a new error dialog
func newDialog(app *tview.Application, theme *Theme) *dialog {
	modal := tview.NewModal()
	errorDialog := &dialog{
		Modal:        modal,
		app:          app,
		previousPage: nil,
	}

	// Initial styling - will be properly themed later
	modal.SetBorder(true).
		SetTitleAlign(tview.AlignCenter)

	errorDialog.applyTheme(theme)
	return errorDialog
}

// Show displays the dialog with the given message
// and returns to the previous page when Enter is pressed
func (d *dialog) Show(previousPage tview.Primitive, title, msg string) {
	d.previousPage = previousPage

	d.SetTitle(title)
	d.SetText(msg).
		ClearButtons().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			d.app.SetRoot(d.previousPage, true)
		})
	d.app.SetRoot(d.Modal, true)
}

func (d *dialog) applyTheme(theme *Theme) {
	colors := theme.GetColors()
	d.SetBorderColor(colors.BorderFocus)
	d.SetTextColor(colors.Foreground)
	d.SetBackgroundColor(colors.Background)
	d.SetTitleColor(colors.BorderFocus)

	d.SetButtonBackgroundColor(colors.Button)
	d.SetButtonTextColor(colors.ButtonText)
	d.SetBorderStyle(tcell.StyleDefault.
		Foreground(colors.BorderFocus).
		Background(colors.Background))
	d.SetButtonActivatedStyle(tcell.StyleDefault.
		Background(colors.ButtonFocus).
		Foreground(colors.ButtonTextFocus))
	d.SetButtonStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))
}
