package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// The sidebar displays tree representing either database or local files information.
// At the top of each tree, the database name or the fixed string "local" is displayed.
// The trees show tables associated with the database or files read from the local system.
type sidebar struct {
	tree *tview.TreeView
}

// newSidebar creates a new sidebar.
func newSidebar() *sidebar {
	tree := tview.NewTreeView()

	tree.SetTopLevel(1)
	root := tview.NewTreeNode(".").SetColor(tcell.ColorRed)
	tree.SetTitle("Databases")
	tree.SetCurrentNode(root)
	tree.SetTitleAlign(tview.AlignLeft)
	tree.SetBorder(true)

	rootNode := tview.NewTreeNode("-")
	tree.SetRoot(rootNode)
	tree.SetCurrentNode(rootNode)

	add := func(_ *tview.TreeNode, _ string) {
		// Do nothing.
	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		add(node, "path")
	})

	return &sidebar{
		tree: tree,
	}
}
