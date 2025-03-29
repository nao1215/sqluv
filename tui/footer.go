package tui

import (
	"maps"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// footer represents the footer component showing keyboard shortcuts
type footer struct {
	*tview.InputField
	shortcuts    map[string]string // key=shortcut, value=description
	theme        *Theme
	searchActive bool
}

// newFooter creates a new footer with keyboard shortcuts
func newFooter(theme *Theme) *footer {
	inputField := tview.NewInputField().
		SetLabel("").
		SetText("").
		SetFieldWidth(0).
		SetLabelWidth(0)
	inputField.SetFieldStyle(tcell.StyleDefault.
		Background(theme.GetColors().Background).
		Foreground(theme.GetColors().Foreground))
	inputField.SetLabelStyle(tcell.StyleDefault.
		Background(theme.GetColors().Background).
		Foreground(theme.GetColors().Header))
	inputField.SetFieldBackgroundColor(theme.GetColors().Background)
	inputField.SetBackgroundColor(theme.GetColors().Background)
	inputField.SetBorder(false)
	inputField.SetDisabled(true)

	footer := &footer{
		InputField: inputField,
		shortcuts:  make(map[string]string),
		theme:      theme,
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

// ActivateSearch enables the search field in the footer.
func (f *footer) ActivateSearch() {
	f.searchActive = true
	f.SetDisabled(false)
	f.SetLabel("Search: ")
	f.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyESC {
			f.update()
			f.SetDisabled(true)
			f.searchActive = false
		}
	})
}

// DeactivateSearch disables the search field.
func (f *footer) DeactivateSearch() {
	f.searchActive = false
	f.SetText("")
	f.SetLabel("")
	f.SetDisabled(true)
}

// update updates the footer text with the current shortcuts
func (f *footer) update() {
	text := ""
	for _, key := range slices.Sorted(maps.Keys(f.shortcuts)) {
		if text != "" {
			text += " | "
		}
		// Use the header color from theme for shortcuts
		text += key + ": " + f.shortcuts[key]
	}
	text += strings.Repeat(" ", 100) // workaround for the footer background color
	f.SetLabel(text)
}

// setDefaulShortcut changes the shortcuts to the setDefaulShortcut screen.
func (f *footer) setDefaulShortcut() {
	f.clearShortcuts()
	f.addShortcut("Ctrl-d", "Quit")
	f.addShortcut("TAB,F1-F3", "Change focus")
	f.addShortcut("Ctrl-t", "Theme")
	f.addShortcut("Ctrl-h", "History")
	f.addShortcut("Ctrl-e", "Exec Query")
	f.update()
}

// setSidebarShortcut changes the shortcuts to the sidebar screen.
func (f *footer) setSidebarShortcut() {
	f.clearShortcuts()
	f.addShortcut("/", "Search")
	f.addShortcut("ESC", "Clear search")
	f.update()
}

func (f *footer) applyTheme(theme *Theme) {
	f.theme = theme
	colors := theme.GetColors()
	f.SetFieldStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))
	f.SetLabelStyle(tcell.StyleDefault.
		Background(theme.GetColors().Background).
		Foreground(theme.GetColors().Foreground))
	f.SetPlaceholderStyle(tcell.StyleDefault.
		Background(theme.GetColors().Background).
		Foreground(theme.GetColors().Foreground))

	f.update()
	if f.HasFocus() {
		f.SetBorderColor(colors.BorderFocus)
	} else {
		f.SetBorderColor(colors.Border)
	}
}

func (f *footer) isActiveSearch() bool {
	return f.searchActive
}
