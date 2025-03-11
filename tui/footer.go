package tui

import (
	"github.com/rivo/tview"
)

// footer represents the footer component showing keyboard shortcuts
type footer struct {
	*tview.TextView
	shortcuts map[string]string // key=shortcut, value=description
	theme     *Theme
}

// newFooter creates a new footer with keyboard shortcuts
func newFooter(theme *Theme) *footer {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	// Apply theme colors instead of hardcoded ones
	colors := theme.GetColors()
	textView.SetTextColor(colors.Foreground)
	textView.SetBackgroundColor(colors.Background)

	footer := &footer{
		TextView:  textView,
		shortcuts: make(map[string]string),
		theme:     theme, // Store theme reference
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

	colors := f.theme.GetColors()

	// Build shortcuts text with formatting
	text := ""
	for key, description := range f.shortcuts {
		if text != "" {
			text += " | "
		}
		// Use the header color from theme for shortcuts
		text += "[" + colors.Header.String() + "]" + key +
			"[" + colors.Foreground.String() + "]: " + description
	}
	f.SetText(text)
}

// setDefaulShortcut changes the shortcuts to the setDefaulShortcut screen.
func (f *footer) setDefaulShortcut() {
	f.clearShortcuts()
	f.addShortcut("Ctrl-D", "Quit")
	f.addShortcut("Esc", "Quit")
	f.addShortcut("TAB", "Change focus")
	f.addShortcut("Ctrl-T", "Theme selector")
	f.update()
}

func (f *footer) applyTheme(theme *Theme) {
	f.theme = theme
	colors := theme.GetColors()
	f.SetBackgroundColor(colors.Background)

	f.update()

	if f.HasFocus() {
		f.SetBorderColor(colors.BorderFocus)
	} else {
		f.SetBorderColor(colors.Border)
	}
}
