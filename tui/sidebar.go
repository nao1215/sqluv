package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/rivo/tview"
)

// The sidebar displays tree representing either database or local files information.
// At the top of each tree, the database name or the fixed string "local" is displayed.
// The trees show tables associated with the database or files read from the local system.
type sidebar struct {
	*tview.TreeView
	theme     *Theme
	allTables []*model.Table
	dbName    string
}

// newSidebar creates a new sidebar.
func newSidebar(theme *Theme) *sidebar {
	tree := tview.NewTreeView()
	colors := theme.GetColors()

	tree.SetTopLevel(1)
	root := tview.NewTreeNode(".").SetColor(tcell.ColorRed)
	tree.SetTitle("Databases")
	tree.SetTitleColor(colors.Header)
	tree.SetCurrentNode(root)
	tree.SetTitleAlign(tview.AlignLeft)
	tree.SetBorder(true)

	rootNode := tview.NewTreeNode("-")
	tree.SetRoot(rootNode)
	tree.SetCurrentNode(rootNode)

	tree.SetBackgroundColor(colors.Background)
	tree.SetBorderColor(colors.Border)
	tree.SetGraphicsColor(colors.Foreground)

	sb := &sidebar{
		TreeView: tree,
		theme:    theme,
	}
	sb.applyTheme(theme)
	return sb
}

// updateTables updates the sidebar with tables and, on table click, shows the column names.
func (s *sidebar) update(tables []*model.Table, dbName string) {
	s.allTables = tables
	s.dbName = dbName
	s.updateTables(tables, dbName)
}

func (s *sidebar) updateTables(tables []*model.Table, dbName string) {
	root := s.GetRoot()
	if root == nil {
		root = tview.NewTreeNode("Databases")
		s.SetRoot(root)
	}
	root.ClearChildren()

	colors := s.theme.GetColors()
	root.SetTextStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))

	// Create the database node.
	dbNode := tview.NewTreeNode(dbName).
		SetSelectable(true).
		SetTextStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.Foreground))
	root.AddChild(dbNode)

	seen := make(map[string]bool)
	// Add tables under the database node.
	for _, table := range tables {
		// Skip duplicate table names.
		if seen[table.Name()] {
			continue
		}
		seen[table.Name()] = true

		// Capture table in local variable for closure.
		tbl := table
		tableNode := tview.NewTreeNode("▷ " + tbl.Name()).
			SetSelectable(true).
			SetReference(tbl).
			SetTextStyle(tcell.StyleDefault.
				Background(colors.Background).
				Foreground(colors.Foreground))

		// When a table node is selected, toggle its column names.
		tableNode.SetSelectedFunc(func() {
			if len(tableNode.GetChildren()) == 0 {
				for _, col := range tbl.Header() {
					colNode := tview.NewTreeNode(col).
						SetSelectable(false).
						SetTextStyle(tcell.StyleDefault.
							Background(colors.Background).
							Foreground(colors.Foreground))
					tableNode.AddChild(colNode)
				}
			} else {
				tableNode.ClearChildren()
			}
			applyThemeToTreeNodes(tableNode, colors)
			s.SetCurrentNode(tableNode)
		})
		dbNode.AddChild(tableNode)
	}
	s.SetCurrentNode(root)
	applyThemeToTreeNodes(root, colors)
}

func (s *sidebar) applyTheme(theme *Theme) {
	colors := theme.GetColors()

	s.SetBackgroundColor(colors.Background)
	s.SetTitleColor(colors.Header)

	// Apply border color based on focus state
	if s.HasFocus() {
		s.SetBorderColor(colors.BorderFocus)
	} else {
		s.SetBorderColor(colors.Border)
	}

	s.SetGraphicsColor(colors.Foreground)

	// Apply theme to all existing tree nodes
	root := s.GetRoot()
	if root != nil {
		applyThemeToTreeNodes(root, colors)
	}
}

// Helper function to recursively apply theme to all tree nodes
func applyThemeToTreeNodes(node *tview.TreeNode, colors ThemeColors) {
	// Don't override specific colors set for special nodes (like database names)
	if node.GetColor() == tcell.ColorDefault {
		node.SetColor(colors.Foreground)
	}

	// Set the background color for the node (this was missing)
	node.SetTextStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))

	// Apply to all children recursively
	for _, child := range node.GetChildren() {
		applyThemeToTreeNodes(child, colors)
	}
}

// filterTables filters the sidebar table nodes by query using fuzzy matching.
func (s *sidebar) filterTables(query string) {
	// Use fuzzy matching to select tables.
	matches := []*model.Table{}
	for _, table := range s.allTables {
		if fuzzy.Match(query, table.Name()) {
			matches = append(matches, table)
		}
	}
	s.updateTables(matches, s.dbName)
}
