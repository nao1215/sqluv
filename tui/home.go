package tui

import "github.com/rivo/tview"

// home represents the home window.
type home struct {
	flex *tview.Flex
}

// newHome creates a new home window.
func newHome() *home {
	sidebar := newSidebar().tree
	textArea := newQueryTextArea().textArea
	resultTable := newQueryResultTable().table

	// Create a flex for the query input and results (vertical layout)
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textArea, 6, 0, false).
		AddItem(resultTable, 0, 1, false)

	// Main layout with sidebar and right panel (horizontal layout)
	flex := tview.NewFlex().
		AddItem(sidebar, 0, 1, false).
		AddItem(rightPanel, 0, 5, false)

	return &home{
		flex: flex,
	}
}
