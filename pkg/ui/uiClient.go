package ui

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/thebsdbox/kuiz/pkg/quiz/types"
)

//ClientView is the main user interface
type ClientView struct {
	App           *tview.Application
	questions     *tview.TreeView
	leftFlexPanel *tview.Flex

	// Main UI components including the question box and the flex component for answers
	// Questions UI
	question *tview.TextView
	answers  *tview.Flex

	// score UI
	score *tview.TextView

	// welcome UI
	welcome *tview.TextView

	// Answered UI
	answered *tview.TextView
}

// BuildClientView -
func BuildClientView() (*ClientView, error) {
	v := ClientView{
		App: tview.NewApplication(),

		// Left Panel
		leftFlexPanel: tview.NewFlex(),

		// Question UI init
		question: tview.NewTextView(),
		answers:  tview.NewFlex(),

		// Score UI init
		score: tview.NewTextView(),

		// Welcome UI init
		welcome: tview.NewTextView(),
		// Right panel (list of asked questions)
		questions: tview.NewTreeView(),

		// Answered UI
		answered: tview.NewTextView(),
	}

	v.addQuestionListView()
	v.addQuestionTextView()

	v.leftFlexPanel.AddItem(tview.NewFlex().SetDirection(tview.FlexRow), 0, 2, false).AddItem(v.welcome, 0, 0, false)
	v.leftFlexPanel.AddItem(v.questions, 30, 1, false)

	v.App.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})
	v.App.SetRoot(v.leftFlexPanel, true)
	v.App.SetFocus(v.answers)

	return &v, nil
}

// SetFlexPanelToWelcome - changes the UI to display the question
func (v *ClientView) SetFlexPanelToWelcome(welcomeText string) {
	// Remove things on the panel
	v.App.QueueUpdateDraw(func() {
		v.leftFlexPanel.Clear()
		// Add the Quesion panel
		v.welcome.SetText(fmt.Sprintf("\n\n\n\n\nThankyou for joining this kuiz! \n\n\nTodays topic:\n\n\n\n\n\n\"%s\" ", welcomeText))
		//v.welcome.SetBackgroundColor(tcell.ColorMediumSeaGreen)
		v.welcome.SetBorder(true)
		v.welcome.SetTextAlign(tview.AlignCenter)
		v.leftFlexPanel.AddItem(v.welcome, 0, 1, true)
	})
}

// SetFlexPanelToQuestion - changes the UI to display the question
func (v *ClientView) SetFlexPanelToQuestion(question string) {
	// Remove things on the panel
	v.App.QueueUpdateDraw(func() {
		v.leftFlexPanel.Clear()
		// Add the Quesion panel
		questionView := tview.NewTextView().SetText(fmt.Sprintf("\n\n\n\"%s\"\n\n\n Answers revealed soon ... ", question))
		//questionView.SetBackgroundColor(tcell.ColorGreen)
		questionView.SetBorder(true)
		questionView.SetTextAlign(tview.AlignCenter)
		v.leftFlexPanel.AddItem(questionView, 0, 1, true)
	})
}

// SetFlexPanelToAnswers - changes the UI to display the questions
func (v *ClientView) SetFlexPanelToAnswers() {
	// Remove things on the panel
	v.App.QueueUpdateDraw(func() {
		v.leftFlexPanel.Clear()
		// Add the Quesion panel
		v.leftFlexPanel.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(v.question, 4, 1, false).
			AddItem(v.answers, 0, 3, false), 0, 2, false)
		v.leftFlexPanel.AddItem(v.questions, 30, 1, false)
	})
}

// SetFlexPanelToAnswered - changes the UI to display the questions
func (v *ClientView) SetFlexPanelToAnswered() {
	// Remove things on the panel
	v.App.QueueUpdateDraw(func() {
		v.leftFlexPanel.Clear()
		// Add the Answered panel
		v.answered.SetText("\n\n\n\n\n\n\n\n\n\nAnswered...")
		//v.answered.SetBackgroundColor(tcell.ColorBisque)
		v.answered.SetBorder(true)
		v.answered.SetTextAlign(tview.AlignCenter)
		v.leftFlexPanel.AddItem(v.answered, 0, 1, false)

	})
}

// SetFlexPanelToScore changes the UI to display the Score screen
func (v *ClientView) SetFlexPanelToScore(scoreboard *types.ScoreBoard) {
	scoretext := "\n\n\nTHE SCORES\n\n\n"
	for x := range scoreboard.Users {
		scoretext += fmt.Sprintf("-> [red] %s: [white] %s with a score of -> [green]%d[white]\n", humanize.Ordinal(x+1), scoreboard.Users[x].User.Username, scoreboard.Users[x].Score)
	}

	scoretext += fmt.Sprintf("\n\n\n\n\n\nFastest player -> [red] %s [white]with reflexes to answer a question in [red]%dms[white]", scoreboard.FastestUser, scoreboard.FastestAnswer)
	// Remove things on the panel
	v.App.QueueUpdateDraw(func() {
		v.leftFlexPanel.Clear()
		// text, err := json.MarshalIndent(scoreboard, "", "\t")
		// if err != nil {
		// }
		v.score.SetText(scoretext)
		v.score.SetDynamicColors(true)
		//v.score.SetBackgroundColor(tcell.ColorDarkBlue)
		v.score.SetBorder(true)
		v.score.SetTextAlign(tview.AlignCenter)
		v.leftFlexPanel.AddItem(v.score, 0, 1, true)
	})
}

// UpdateQuestionText - Updates the list of past questions
func (v *ClientView) UpdateQuestionText(q string) {
	v.App.QueueUpdateDraw(func() {
		// Update Question Box
		v.question.SetText(q)
		// Add question to the log
		v.questions.GetRoot().AddChild(tview.NewTreeNode(q))
	})
}
