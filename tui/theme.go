package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/config"
	"github.com/rivo/tview"
)

// Theme manages the color scheme for the application
type Theme struct {
	color *config.ColorConfig
	app   *tview.Application
}

// NewTheme creates a new theme manager
func NewTheme(color *config.ColorConfig, app *tview.Application) *Theme {
	return &Theme{
		color: color,
		app:   app,
	}
}

// GetColors returns the current colors based on the active color scheme
func (t *Theme) GetColors() ThemeColors {
	scheme := t.color.CurrentScheme

	return ThemeColors{
		Background:      config.GetTcellColor(scheme.Background),
		Foreground:      config.GetTcellColor(scheme.Foreground),
		Border:          config.GetTcellColor(scheme.Border),
		BorderFocus:     config.GetTcellColor(scheme.BorderFocus),
		Selection:       config.GetTcellColor(scheme.Selection),
		SelectionText:   config.GetTcellColor(scheme.SelectionText),
		Header:          config.GetTcellColor(scheme.Header),
		Button:          config.GetTcellColor(scheme.Button),
		ButtonFocus:     config.GetTcellColor(scheme.ButtonFocus),
		ButtonText:      config.GetTcellColor(scheme.ButtonText),
		ButtonTextFocus: config.GetTcellColor(scheme.ButtonTextFocus),
	}
}

// ThemeColors holds the actual tcell.Color values for the application
type ThemeColors struct {
	Background      tcell.Color
	Foreground      tcell.Color
	Border          tcell.Color
	BorderFocus     tcell.Color
	Selection       tcell.Color
	SelectionText   tcell.Color
	Header          tcell.Color
	Button          tcell.Color
	ButtonFocus     tcell.Color
	ButtonText      tcell.Color
	ButtonTextFocus tcell.Color
}

// ShowColorSchemeSelector displays a modal for selecting the color scheme
// with buttons arranged in rows of 4 in the center of the screen
func (t *Theme) ShowColorSchemeSelector(onClose func()) {
	schemes := t.color.GetSchemeNames()

	// Create a flex container for the buttons with padding
	content := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Add some padding at the top
	content.AddItem(nil, 2, 0, false)

	// Calculate how many rows we need
	rowCount := (len(schemes) + 3) / 4 // Ceiling division by 4

	// Keep track of all buttons to set up navigation
	var allButtons []*tview.Button
	var rowButtons [][]*tview.Button

	// Create each row with up to 4 buttons
	for i := 0; i < rowCount; i++ {
		// Create a flex for this row with padding on sides
		row := tview.NewFlex().SetDirection(tview.FlexColumn)

		// Add left padding
		row.AddItem(nil, 0, 1, false)

		// Track buttons in this row
		buttonsInRow := []*tview.Button{}

		// Add up to 4 buttons per row
		for j := 0; j < 4; j++ {
			index := i*4 + j
			if index >= len(schemes) {
				// Skip if we've processed all schemes
				continue
			}

			schemeName := schemes[index]
			button := tview.NewButton(schemeName)

			// Get colors for styling
			colors := t.GetColors()
			button.SetStyle(tcell.StyleDefault.
				Background(colors.Button).
				Foreground(colors.ButtonText))
			button.SetActivatedStyle(tcell.StyleDefault.
				Background(colors.ButtonFocus).
				Foreground(colors.ButtonTextFocus))
			button.SetBorderStyle(tcell.StyleDefault.
				Background(colors.Button).
				Foreground(colors.ButtonText))
			button.SetDisabledStyle(tcell.StyleDefault.
				Background(colors.Button).
				Foreground(colors.ButtonText))
			button.SetLabelColor(colors.ButtonText)

			// Use closure to capture the correct scheme name
			func(name string) {
				button.SetSelectedFunc(func() {
					t.color.SetScheme(name)
					if onClose != nil {
						onClose()
					}
				})
			}(schemeName)

			row.AddItem(button, 0, 1, true)
			allButtons = append(allButtons, button)
			buttonsInRow = append(buttonsInRow, button)
		}

		rowButtons = append(rowButtons, buttonsInRow)

		// Add right padding (flexible to center the buttons)
		row.AddItem(nil, 0, 1, false)

		// Add the row to the content flex with fixed height
		content.AddItem(row, 3, 0, true)
	}

	// Add a cancel button row at the bottom
	cancelRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Add left padding for centering
	cancelRow.AddItem(nil, 0, 1, false)

	cancelButton := tview.NewButton("Cancel")

	// Style the cancel button
	colors := t.GetColors()
	cancelButton.SetStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))
	cancelButton.SetActivatedStyle(tcell.StyleDefault.
		Background(colors.ButtonFocus).
		Foreground(colors.ButtonTextFocus))
	cancelButton.SetBorderStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))
	cancelButton.SetDisabledStyle(tcell.StyleDefault.
		Background(colors.Button).
		Foreground(colors.ButtonText))
	cancelButton.SetLabelColor(colors.ButtonText)

	cancelButton.SetSelectedFunc(func() {
		if onClose != nil {
			onClose()
		}
	})

	cancelRow.AddItem(cancelButton, 10, 0, true)
	allButtons = append(allButtons, cancelButton)

	// Add right padding for centering
	cancelRow.AddItem(nil, 0, 1, false)

	// Add spacing between color buttons and cancel button
	content.AddItem(nil, 1, 0, false)
	content.AddItem(cancelRow, 3, 0, true)
	content.AddItem(nil, 2, 0, false) // Bottom padding

	// Set border around the whole content
	frame := tview.NewFrame(content).
		SetBorders(1, 1, 1, 1, 1, 1).
		AddText("Select Color Scheme", true, tview.AlignCenter, colors.Header)

	frame.SetBorderColor(colors.BorderFocus)
	frame.SetBackgroundColor(colors.Background)

	// Create pages to show the modal with a dimmed background
	pages := tview.NewPages().
		AddPage("background", tview.NewBox().SetBackgroundColor(tcell.ColorBlack), true, true).
		AddPage("modal", frame, true, true)

	// Configure keyboard navigation between buttons
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Get the currently focused primitive
		focused := t.app.GetFocus()

		// Find the current button index
		currentIndex := -1
		for i, btn := range allButtons {
			if focused == btn {
				currentIndex = i
				break
			}
		}

		// If no button is focused, focus the first one
		if currentIndex == -1 {
			t.app.SetFocus(allButtons[0])
			return nil
		}

		// Find the current row and column
		currentRow := -1
		currentCol := -1

		for i, row := range rowButtons {
			for j, btn := range row {
				if focused == btn {
					currentRow = i
					currentCol = j
					break
				}
			}
			if currentRow != -1 {
				break
			}
		}

		switch event.Key() {
		case tcell.KeyTab:
			// Move to the next button
			nextIndex := (currentIndex + 1) % len(allButtons)
			t.app.SetFocus(allButtons[nextIndex])
			return nil

		case tcell.KeyBacktab:
			// Move to the previous button
			nextIndex := (currentIndex - 1 + len(allButtons)) % len(allButtons)
			t.app.SetFocus(allButtons[nextIndex])
			return nil

		case tcell.KeyRight:
			// If we're in a row (not the cancel button)
			if currentRow != -1 && currentCol != -1 {
				// If not at the end of the row
				if currentCol < len(rowButtons[currentRow])-1 {
					t.app.SetFocus(rowButtons[currentRow][currentCol+1])
					return nil
				} else if currentRow < len(rowButtons)-1 {
					// Move to the first button of the next row
					t.app.SetFocus(rowButtons[currentRow+1][0])
					return nil
				}
			}

		case tcell.KeyLeft:
			// If we're in a row (not the cancel button)
			if currentRow != -1 && currentCol != -1 {
				// If not at the beginning of the row
				if currentCol > 0 {
					t.app.SetFocus(rowButtons[currentRow][currentCol-1])
					return nil
				} else if currentRow > 0 {
					// Move to the last button of the previous row
					lastCol := len(rowButtons[currentRow-1]) - 1
					t.app.SetFocus(rowButtons[currentRow-1][lastCol])
					return nil
				}
			}

		case tcell.KeyDown:
			// If we're in a row and not the last row
			if currentRow != -1 && currentRow < len(rowButtons)-1 {
				// Try to move to the same column in the next row, or the last column if that row is shorter
				nextCol := currentCol
				if nextCol >= len(rowButtons[currentRow+1]) {
					nextCol = len(rowButtons[currentRow+1]) - 1
				}
				t.app.SetFocus(rowButtons[currentRow+1][nextCol])
				return nil
			} else if currentRow != -1 {
				// Move from the last row to the cancel button
				t.app.SetFocus(cancelButton)
				return nil
			}

		case tcell.KeyUp:
			// If on the cancel button, move to the last row
			if focused == cancelButton && len(rowButtons) > 0 {
				lastRow := len(rowButtons) - 1
				// Focus the middle button in the last row
				midCol := len(rowButtons[lastRow]) / 2
				t.app.SetFocus(rowButtons[lastRow][midCol])
				return nil
			} else if currentRow > 0 {
				// Move to the same column in the previous row
				nextCol := currentCol
				if nextCol >= len(rowButtons[currentRow-1]) {
					nextCol = len(rowButtons[currentRow-1]) - 1
				}
				t.app.SetFocus(rowButtons[currentRow-1][nextCol])
				return nil
			}
		}
		return event
	})

	t.app.SetRoot(pages, true)

	// Set initial focus on the first button
	if len(allButtons) > 0 {
		t.app.SetFocus(allButtons[0])
	}
}

// ApplyButtonTheme apply theme to a button
func (t *Theme) ApplyButtonTheme(button *tview.Button, focused bool) {
	colors := t.GetColors()

	if focused {
		button.SetLabelColor(colors.ButtonTextFocus)
		button.SetBackgroundColor(colors.ButtonFocus)
	} else {
		button.SetLabelColor(colors.ButtonText)
		button.SetBackgroundColor(colors.Button)
	}
}

// ApplyTextViewTheme apply theme to a text view
func (t *Theme) ApplyTextViewTheme(textView *tview.TextView, focused bool) {
	colors := t.GetColors()

	textView.SetTextColor(colors.Foreground)
	textView.SetBackgroundColor(colors.Background)

	if focused {
		textView.SetBorderColor(colors.BorderFocus)
	} else {
		textView.SetBorderColor(colors.Border)
	}
}

// ApplyTableTheme apply theme to a table
func (t *Theme) ApplyTableTheme(table *tview.Table, focused bool) {
	colors := t.GetColors()

	table.SetBackgroundColor(colors.Background)

	if focused {
		table.SetBorderColor(colors.BorderFocus)
		table.SetSelectedStyle(tcell.StyleDefault.Background(colors.Selection).Foreground(colors.SelectionText))
	} else {
		table.SetBorderColor(colors.Border)
		table.SetSelectedStyle(tcell.StyleDefault.Background(colors.Background).Foreground(colors.Foreground))
	}
}

// ApplyTextAreaTheme apply theme to a text area
func (t *Theme) ApplyTextAreaTheme(textArea *tview.TextArea, focused bool) {
	colors := t.GetColors()

	textArea.SetTitleColor(colors.Foreground)
	textArea.SetBackgroundColor(colors.Background)

	if focused {
		textArea.SetBorderColor(colors.BorderFocus)
	} else {
		textArea.SetBorderColor(colors.Border)
	}
}
