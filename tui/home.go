package tui

import "github.com/rivo/tview"

// home represents the home window.
type home struct {
	flex          *tview.Flex
	sidebar       *sidebar
	queryTextArea *queryTextArea
	executeButton *executeButton
	historyButton *historyButton
	resultTable   *queryResultTable
	footer        *footer
	dialog        *dialog
	rowStatistics *rowStatistics
}

// newHome creates a new home window.
func newHome(app *tview.Application, theme *Theme) *home {
	sidebarComponent := newSidebar(theme)
	textArea := newQueryTextArea(theme)
	executeButton := newExecuteButton(theme)
	historyButton := newHistoryButton(theme)
	resultTableComponent := newQueryResultTable(theme)
	footerComponent := newFooter(theme)
	rowStatisticsComponent := newRowStatistics(theme)

	buttonPanel := tview.NewFlex().
		AddItem(executeButton, 0, 1, false).
		AddItem(historyButton, 0, 1, false)

	// Create a flex for the query input and results (vertical layout)
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textArea, 6, 0, true).
		AddItem(buttonPanel, 1, 0, false).
		AddItem(resultTableComponent, 0, 1, false).
		AddItem(rowStatisticsComponent, 1, 0, false)

	// Main content with sidebar and right panel (horizontal layout)
	mainContent := tview.NewFlex().
		AddItem(sidebarComponent, 0, 1, false).
		AddItem(rightPanel, 0, 5, false)

	// Create the main layout with content at top and footer at bottom
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainContent, 0, 1, false).
		AddItem(footerComponent, 1, 0, false)

	return &home{
		flex:          mainFlex,
		sidebar:       sidebarComponent,
		queryTextArea: textArea,
		executeButton: executeButton,
		historyButton: historyButton,
		resultTable:   resultTableComponent,
		dialog:        newDialog(app, theme),
		footer:        footerComponent,
		rowStatistics: rowStatisticsComponent,
	}
}

// applyTheme applies the current theme to all components
func (h *home) applyTheme(theme *Theme) {
	// Apply theme to the flex container itself
	colors := theme.GetColors()
	h.flex.SetBackgroundColor(colors.Background)

	// Update each component with the theme
	h.sidebar.applyTheme(theme)
	h.queryTextArea.applyTheme(theme)
	h.executeButton.applyTheme(theme)
	h.historyButton.applyTheme(theme)
	h.resultTable.applyTheme(theme)
	h.footer.applyTheme(theme)
	h.rowStatistics.applyTheme(theme)
	h.dialog.applyTheme(theme)
}
