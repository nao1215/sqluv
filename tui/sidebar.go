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

	return &sidebar{
		TreeView: tree,
	}
}

// Add this method to your sidebar struct
func (s *sidebar) addLocalFiles(files []*model.File) {
	// Find or create the "local" root node
	root := s.GetRoot()
	var localNode *tview.TreeNode

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
			SetExpanded(true)
		root.AddChild(localNode)
	} else {
		localNode.ClearChildren()
	}

	for _, file := range files {
		name := file.NameWithoutExt()
		fileNode := tview.NewTreeNode(name).
			SetSelectable(true).
			SetReference(file) // Store the file reference for later use
		localNode.AddChild(fileNode)
	}
}
