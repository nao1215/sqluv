package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// footer represents the footer component showing keyboard shortcuts
type footer struct {
	*tview.TextView
	shortcuts map[string]string // key=shortcut, value=description
}

// newFooter creates a new footer with keyboard shortcuts
func newFooter() *footer {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorDefault)

	footer := &footer{
		TextView:  textView,
		shortcuts: make(map[string]string),
	}
	return footer
}

// addShortcut adds a keyboard shortcut to the footer
func (f *footer) addShortcut(key, description string) {
	f.shortcuts[key] = description
	f.update()
}

// clearShortcuts removes all shortcuts
func (f *footer) clearShortcuts() {
	f.shortcuts = make(map[string]string)
	f.update()
}

// update updates the footer text with the current shortcuts
func (f *footer) update() {
	f.Clear()

	// Build shortcuts text with formatting
	text := ""
	for key, description := range f.shortcuts {
		if text != "" {
			text += " | "
		}
		text += "[yellow]" + key + "[white]: " + description
	}
	f.SetText(text)
}

// setDefaulShortcut changes the shortcuts to the setDefaulShortcut screen.
func (f *footer) setDefaulShortcut() {
	f.clearShortcuts()
	f.addShortcut("Ctrl-D", "Quit")
	f.addShortcut("Esc", "Quit")
	f.addShortcut("TAB", "Change focus")
	f.update()
}
