package tui

import "github.com/rivo/tview"

// home represents the home window.
type home struct {
	flex          *tview.Flex
	sidebar       *sidebar
	queryTextArea *queryTextArea
	executeButton *executeButton
	resultTable   *queryResultTable
	footer        *footer
	errorDialog   *errorDialog
}

// newHome creates a new home window.
func newHome(app *tview.Application) *home {
	sidebarComponent := newSidebar()
	textArea := newQueryTextArea()
	executeButton := newExecuteButton()
	resultTableComponent := newQueryResultTable()
	footerComponent := newFooter()

	// Create a flex for the query input and results (vertical layout)
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textArea, 6, 0, true).
		AddItem(executeButton, 1, 0, false).
		AddItem(resultTableComponent, 0, 1, false)

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
		resultTable:   resultTableComponent,
		errorDialog:   newErrorDialog(app),
		footer:        footerComponent,
	}
}
