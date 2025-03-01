package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// errorDialog represents a modal dialog for displaying error messages
type errorDialog struct {
	*tview.Modal
	previousPage tview.Primitive
	app          *tview.Application
}

// newErrorDialog creates a new error dialog
func newErrorDialog(app *tview.Application) *errorDialog {
	modal := tview.NewModal()
	modal.SetBorder(true).
		SetBorderColor(tcell.ColorRed).
		SetTitleColor(tcell.ColorRed).
		SetTitleAlign(tview.AlignCenter).
		SetTitle(" ERROR ")

	return &errorDialog{
		Modal:        modal,
		app:          app,
		previousPage: nil,
	}
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
