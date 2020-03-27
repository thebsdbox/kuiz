package server

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
	"github.com/thebsdbox/kuiz/pkg/quiz/types"
	"github.com/thebsdbox/kuiz/pkg/ui"
)

//var addr = flag.String("addr", "localhost:8080", "http service address")

// The server package is the "engine" for sending questions to the clients

// Manager - handles all required objects for the client to work
type Manager struct {
	addr string
	v    *ui.ServerView
	q    *types.Quiz
	stop chan bool
	qm   *questionManager
	s    types.ScoreBoard
}

//NewServer - generates a new server object
func NewServer(address string, quiz *types.Quiz, v *ui.ServerView) *Manager {
	m := Manager{
		addr: address,
		v:    v,
		q:    quiz,
		stop: make(chan bool, 1),
	}

	// Set the fastest answer to something big
	m.s.FastestAnswer = time.Hour.Milliseconds()

	// TODO - strip out this function
	v.AddQuestionHandler(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}

		//New Question manager for this question
		m.qm = &questionManager{}

		question := reference.(types.Question)
		m.qm.question = question

		// // Strip answer
		// switch question.Type {
		// case types.QUIZMULTI:
		// 	m.qm.correctAnswer = question.MultiChoiceQuestion.Answer
		// 	//question.MultiChoiceQuestion.Answer = ""

		// case types.QUIZYESNO:
		// 	m.qm.correctAnswer = question.YesNoChoiceQuestion.Answer

		// 	//question.YesNoChoiceQuestion.Answer = ""
		// }
		// Wrap the data to send to clients
		wrappedData := types.EncodeData(types.TYPEQUIZ, question)
		for x := range clients {
			// send to all clients (ignore the dead or failed clients)
			if clients[x] != nil {
				go clients[x].WriteJSON(wrappedData)
			}
		}

		// Time out for question at this point, in a seperate go routine to not block the UI
		go func() {
			time.Sleep(time.Second * 10)
			m.processAnswers()
			m.postScore()
		}()
	})

	// Allows the server operator to bring up the scoreboard at any point
	v.AddScoreHandler(m.postScore)

	return &m
}

//StartServer will start the kuiz server
func (m *Manager) StartServer() error {

	// Set the handler for the web server
	http.HandleFunc("/", m.server)

	// TODO - remove this
	m.v.UpdateQuestions(m.q)
	// Start the server
	return http.ListenAndServe(m.addr, nil)
}

// StopServer will close all connections and wait one second before ending
func (m *Manager) StopServer() {
	// Send close message
	for x := range clients {
		// send to all clients (ignore the dead or failed clients)
		go clients[x].WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}
	m.v.App.Stop()
	// Wait for all things to be cleared up
	<-time.After(time.Second)
}
