package ui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var activeButtons []tview.Button
var activeIndex int

func (v *ClientView) addQuestionListView() {
	rootDir := "Questions Log"
	root := tview.NewTreeNode(rootDir).SetColor(tcell.ColorGreen)
	v.questions.SetRoot(root).SetCurrentNode(root).SetTopLevel(0)
	v.questions.SetBorder(true)
	v.questions.SetBackgroundColor(tcell.ColorDefault)
}

// Sets up the question Text view box
func (v *ClientView) addQuestionTextView() {
	v.question.SetBorder(true)
	v.question.SetTitle("Question")
}

// UpdateAnswers will update the answers flex with all answer choices
func (v *ClientView) UpdateAnswers(timeToReveal time.Duration, answers []string) error {

	// Draw answers
	v.App.QueueUpdateDraw(func() {
		if len(answers) == 0 {
			//return fmt.Errorf("No answers specified")
		}
		//Clear the existing form
		v.answers.Clear()
		// Set the border and the title
		v.answers.SetBorder(true)
		v.answers.SetTitle("Answers")

		// TODO - add more colours?
		colours := []tcell.Color{tcell.ColorRed, tcell.ColorGreen, tcell.ColorBlue, tcell.ColorOrange}

		// Clean the array of active buttons
		activeButtons = nil
		// reset the active button
		activeIndex = 0

		// Build our first row for answers
		AnswerRow := tview.NewFlex().SetDirection(tview.FlexRow)

		// Iterate over the answers
		for x := range answers {
			// Set teh callback function
			newButton := tview.NewButton(answers[x]).SetSelectedFunc(func() {
				b := v.App.GetFocus().(*tview.Button)
				answerText := b.GetLabel()
				answerHandler(answerText)
				//v.SetFlexPanelToAnswered()
				v.leftFlexPanel.Clear()
			})
			// Add our new button to the list of buttons for "tabbing" over
			activeButtons = append(activeButtons, *newButton)
			// Set colour and button
			newButton.SetBackgroundColor(colours[x])
			newButton.SetBorder(true)
			// Add capture for the button
			newButton.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
				if key.Name() == "Tab" {

					v.App.SetFocus(&activeButtons[activeIndex])
					activeIndex++
					if activeIndex == len(activeButtons) {
						activeIndex = 0
					}
					return nil
				}
				return key

			})
			// Add the button to the row
			AnswerRow.AddItem((newButton), 0, 1, true)
			// If we've moved to the second entry and there are only more than two answers, add a new row
			if x == 1 && len(answers) != 2 {
				v.answers.AddItem((AnswerRow), 0, 2, true)
				AnswerRow = tview.NewFlex().SetDirection(tview.FlexRow)
			}
		}
		v.answers.AddItem((AnswerRow), 0, 2, true)
	})
	v.App.SetFocus(&activeButtons[0])
	return nil
}
