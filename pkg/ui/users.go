package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Keeps a list of clients
var clients []client

type client struct {
	name  string
	close chan struct{}
}

//AddUserToView - adds to our list
func (v *ServerView) AddUserToView(c string, close chan struct{}) {
	newClient := client{
		name:  c,
		close: close,
	}

	clients = append(clients, newClient)
	v.syncTree()
}

//DelUserFromView - removes a client
func (v *ServerView) DelUserFromView(c string) {
	for i, v := range clients {
		if v.name == c {
			clients = append(clients[:i], clients[i+1:]...)
		}
	}
	v.syncTree()
}

func (v *ServerView) syncTree() {
	v.App.QueueUpdateDraw(func() {
		root := v.clients.GetRoot()
		root.ClearChildren()
		for x := range clients {
			root.AddChild(tview.NewTreeNode(clients[x].name))
		}
	})
}

func (v *ServerView) addClientView() {
	rootDir := "Users->"
	root := tview.NewTreeNode(rootDir).SetColor(tcell.ColorGreen)
	v.clients.SetRoot(root).SetCurrentNode(root).SetTopLevel(0)
	v.clients.SetBorder(true)
	v.clients.SetBackgroundColor(tcell.ColorDefault)

	v.clients.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			return event

		} else if event.Key() == tcell.KeyDelete {
			node := v.clients.GetCurrentNode()
			v.DelUserFromView(node.GetText())
			return nil

		} else {
			return event
		}
	})

	v.clients.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			node.AddChild(tview.NewTreeNode(path))
			//add(node, path, layerFileSystem)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	v.clients.GetDrawFunc()
}
