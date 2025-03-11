package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/rivo/tview"
)

// The sidebar displays tree representing either database or local files information.
// At the top of each tree, the database name or the fixed string "local" is displayed.
// The trees show tables associated with the database or files read from the local system.
type sidebar struct {
	*tview.TreeView
	theme *Theme
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

// Add this method to your sidebar struct
func (s *sidebar) updateLocalFiles(files []*model.File) {
	// Find or create the "local" root node
	root := s.GetRoot()
	var localNode *tview.TreeNode

	colors := s.theme.GetColors()
	root.SetTextStyle(tcell.StyleDefault.
		Background(colors.Background).
		Foreground(colors.Foreground))

	// Look for existing local node
	for _, child := range root.GetChildren() {
		if child.GetText() == "local" {
			localNode = child
			break
		}
	}

	// If local node doesn't exist, create it
	if localNode == nil {
		localNode = tview.NewTreeNode("local").
			SetSelectable(true).
			SetExpanded(true).
			SetTextStyle(tcell.StyleDefault.
				Background(colors.Background).
				Foreground(colors.Foreground))
		root.AddChild(localNode)
	} else {
		localNode.ClearChildren()
	}

	for _, file := range files {
		name := file.NameWithoutExt()
		fileNode := tview.NewTreeNode(name).
			SetSelectable(true).
			SetReference(file).
			SetTextStyle(tcell.StyleDefault.
				Background(colors.Background).
				Foreground(colors.Foreground))
		localNode.AddChild(fileNode)
	}
	applyThemeToTreeNodes(root, colors)
}

// updateTables update tables in the sidebar
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

	// Create database node
	dbNode := tview.NewTreeNode(dbName).
		SetSelectable(true).
		SetTextStyle(tcell.StyleDefault.
			Background(colors.Background).
			Foreground(colors.Foreground))
	root.AddChild(dbNode)

	// Add tables under the database node
	for _, table := range tables {
		tableNode := tview.NewTreeNode(table.Name()).
			SetSelectable(true).
			SetReference(table).
			SetTextStyle(tcell.StyleDefault.
				Background(colors.Background).
				Foreground(colors.Foreground))
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
