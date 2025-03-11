package tui

import (
	"github.com/rivo/tview"
)

// errorDialog represents a modal dialog for displaying error messages
type errorDialog struct {
	*tview.Modal
	previousPage tview.Primitive
	app          *tview.Application
}

// newErrorDialog creates a new error dialog
func newErrorDialog(app *tview.Application, theme *Theme) *errorDialog {
	modal := tview.NewModal()
	errorDialog := &errorDialog{
		Modal:        modal,
		app:          app,
		previousPage: nil,
	}

	// Initial styling - will be properly themed later
	modal.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle(" ERROR ")

	errorDialog.applyTheme(theme)
	return errorDialog
}

// Show displays the error dialog with the given message
// and returns to the previous page when Enter is pressed
func (d *errorDialog) Show(previousPage tview.Primitive, errorMsg string) {
	d.previousPage = previousPage

	d.SetText(errorMsg).
		ClearButtons().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			d.app.SetRoot(d.previousPage, true)
		})
	d.app.SetRoot(d.Modal, true)
	d.app.SetFocus(d.Modal)
}

func (d *errorDialog) applyTheme(theme *Theme) {
	colors := theme.GetColors()
	d.SetBorderColor(colors.BorderFocus)
	d.SetTextColor(colors.Foreground)
	d.SetBackgroundColor(colors.Background)
	d.SetTitleColor(colors.BorderFocus)

	d.SetButtonBackgroundColor(colors.Button)
	d.SetButtonTextColor(colors.ButtonText)
}
