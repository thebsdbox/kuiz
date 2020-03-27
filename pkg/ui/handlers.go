package ui

import "github.com/rivo/tview"

// This sets up call backs and handlers for UI components to send data back to callers
var answerHandler func(answer string)

// scoreHandler func()
var scoreHandler func()

// ------ Server Handlers ------ //

// AddQuestionHandler Installs teh handler for Questions
func (v *ServerView) AddQuestionHandler(questionHandler func(*tview.TreeNode)) {
	v.questionManager.SetSelectedFunc(questionHandler)
}

//AddScoreHandler - sets the call back when we need to send score details
func (v *ServerView) AddScoreHandler(handler func()) {
	scoreHandler = handler
}

// ------ Client Handlers ------ //

//AddAnswerHandler - sets the call back when an answer needs to be sent
func (v *ClientView) AddAnswerHandler(handler func(answer string)) {
	answerHandler = handler
}
