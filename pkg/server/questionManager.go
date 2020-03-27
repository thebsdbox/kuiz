package server

import (
	"sort"
	"time"

	"github.com/thebsdbox/kuiz/pkg/quiz/types"
)

type questionManager struct {
	question      types.Question
	receipts      []types.QuizReceipt
	answers       []types.QuizAnswer
	fastestAnswer int64
	fastestUser   string
}

func (m *Manager) processAnswers() {
	var correctAnswer string

	question := m.qm.question
	switch question.Type {
	case types.QUIZMULTI:
		correctAnswer = question.MultiChoiceQuestion.Answer
		//question.MultiChoiceQuestion.Answer = ""

	case types.QUIZYESNO:
		correctAnswer = question.YesNoChoiceQuestion.Answer

		//question.YesNoChoiceQuestion.Answer = ""
	}
	// Set the fastest time to something big :-)
	m.qm.fastestAnswer = time.Hour.Milliseconds()

	//->DEBUG
	//fmt.Printf("receipts [%d] / answers [%d]\n", len(m.qm.receipts), len(m.qm.answers))

	// Go through all answers
	for x := range m.qm.answers {
		// If this answer was correct, find when the question was recieved by the client
		//fmt.Printf("%s = %s\n", m.qm.answers[x].Answer, correctAnswer)
		if m.qm.answers[x].Answer == correctAnswer {
			// Find the UID and determine time duration
			for y := range m.qm.receipts {
				if m.qm.answers[x].User.UID == m.qm.receipts[y].User.UID {
					// How long did the answer take?

					// ->DEBUG
					//fmt.Printf("Match for [%s]\n", m.qm.answers[x].User.Username)

					questionDuration := m.qm.answers[x].ReceiptTime.Sub(m.qm.receipts[y].ReceiptTime)
					score := (questionDuration - m.qm.question.TimeToReveal).Milliseconds()
					// ->DEBUG
					//questionDuration1 := m.qm.receipts[x].ReceiptTime.Sub(m.qm.answers[y].ReceiptTime)
					// ->DEBUG

					//fmt.Printf("%d vs %d", questionDuration.Milliseconds(), questionDuration1)

					// Was this the fastest answer
					if m.qm.fastestAnswer > score {
						m.qm.fastestAnswer = score
						m.qm.fastestUser = m.qm.answers[x].User.Username
					}

					// Add to score board
					m.updateScore(&m.qm.answers[x].User, score)
				}
			}
		}
	}
}

// Update the scoreboard
func (m *Manager) updateScore(u *types.User, score int64) {

	// 10000000000
	// 5000000000 --
	// 7783488000 --|_ 2783

	// 10000000000
	// 5000000000 --
	// 5423501000 --|_ 423

	// question time = 10 seconds
	// remove the time we wait, leaving the time to see and answer teh question
	// remove how long the user took

	//fmt.Printf("\n%d\n%d\n%d\n", time.Second*10, m.qm.question.TimeToReveal, questionDuration)
	//score := time.Second*10 -
	// Was this the fastest answer
	if m.s.FastestAnswer > score {
		m.s.FastestAnswer = score
		m.s.FastestUser = u.Username
	}
	// Turn the time into a score
	score = (time.Second * 10).Milliseconds() - score
	var found bool

	for x := range m.s.Users {
		if u.UID == m.s.Users[x].User.UID {
			found = true
			m.s.Users[x].Score += score
		}
	}

	if found == false {
		newScore := types.UserScore{
			User:  u,
			Score: score,
		}
		m.s.Users = append(m.s.Users, newScore)
	}
}

func (m *Manager) postScore() {

	sort.Slice(m.s.Users, func(i, j int) bool {
		return m.s.Users[i].Score > m.s.Users[j].Score
	})
	// Wrap the data to send to clients
	wrappedData := types.EncodeData(types.TYPESCORE, m.s)
	for x := range clients {
		// send to all clients (ignore the dead or failed clients)
		go clients[x].WriteJSON(wrappedData)

	}
}
