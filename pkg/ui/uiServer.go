package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/thebsdbox/kuiz/pkg/quiz/types"
)

// ServerView is the main user interface for the server
type ServerView struct {
	App             *tview.Application
	clients         *tview.TreeView
	flex            *tview.Flex
	questionManager *tview.TreeView
}

// BuildServerView - initialises all parts of the UI
func BuildServerView() (*ServerView, error) {
	v := ServerView{
		App:             tview.NewApplication(),
		clients:         tview.NewTreeView(),
		flex:            tview.NewFlex(),
		questionManager: tview.NewTreeView(),
	}
	v.addClientView()
	v.addQuestionsView()

	v.flex.AddItem(v.clients, 0, 1, false)
	v.flex.AddItem(v.questionManager, 0, 6, true)

	v.App.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})
	v.App.SetRoot(v.flex, true)

	return &v, nil
}

// Sets up the view for the questions
func (v *ServerView) addQuestionsView() {
	rootDir := "Questions:"
	root := tview.NewTreeNode(rootDir).SetColor(tcell.ColorBlue)
	v.questionManager.SetRoot(root).SetCurrentNode(root).SetTopLevel(0)
	v.questionManager.SetBorder(true)
	v.questionManager.SetBackgroundColor(tcell.ColorDefault)
	v.questionManager.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyF1 {
			// Send score screen
			scoreHandler()
			return nil
		}
		return key

	})
}

//UpdateQuestions - adds the questions to the questions view
func (v *ServerView) UpdateQuestions(q *types.Quiz) {
	// Parse the quiz
	v.App.QueueUpdateDraw(func() {
		root := v.questionManager.GetRoot()
		root.ClearChildren()
		root.SetText(q.Name)
		for x := range q.Questions {
			switch q.Questions[x].Type {
			case types.QUIZMULTI:
				newNode := tview.NewTreeNode(q.Questions[x].MultiChoiceQuestion.QuestionText)
				newNode.SetReference(q.Questions[x])
				root.AddChild(newNode)
			case types.QUIZYESNO:
				newNode := tview.NewTreeNode(q.Questions[x].YesNoChoiceQuestion.QuestionText)
				newNode.SetReference(q.Questions[x])
				root.AddChild(newNode)
			}
		}
	})
}
