package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/thebsdbox/kuiz/pkg/quiz/types"
	"github.com/thebsdbox/kuiz/pkg/ui"
)

func (m *ClientManager) processIncoming(data *types.DataWrapper) error {
	switch data.DataType {
	case types.TYPESTATUS:
		var status types.Status
		json.Unmarshal(data.Data, &status)
		if status.Welcome != nil {
			m.v.SetFlexPanelToWelcome(*status.Welcome)
		}
		if status.Error != nil {
			ui.ErrorMessage(fmt.Errorf("%s", *status.Error))
			m.StopClient()
		}

	case types.TYPESCORE:
		//fmt.Printf("%s", data.Data)
		var score types.ScoreBoard
		err := json.Unmarshal(data.Data, &score)
		if err != nil {
			return err
		}
		m.v.SetFlexPanelToScore(&score)
	case types.TYPEUSER:
		// We should recieve a key at this point
		json.Unmarshal(data.Data, &m.u)

	case types.TYPEQUIZ:

		// Let Server know we've recieved the question
		var question types.Question
		err := json.Unmarshal(data.Data, &question)
		if err != nil {
			return err
		}

		// Inform the server that we've recieved the question
		receipt := types.QuizReceipt{
			ReceiptTime: time.Now(),
			User:        *m.u,
		}

		dw := types.EncodeData(types.QUIZRECEIPT, receipt)
		err = m.c.WriteJSON(dw)
		if err != nil {
			return err
		}
		// Parse the question
		switch question.Type {
		case types.QUIZMULTI:
			// Show Question
			m.v.SetFlexPanelToQuestion(question.MultiChoiceQuestion.QuestionText)
			// Wait
			time.Sleep(question.TimeToReveal)
			// Update question text
			m.v.UpdateQuestionText(question.MultiChoiceQuestion.QuestionText)
			m.v.UpdateAnswers(question.TimeToReveal, question.MultiChoiceQuestion.Answers)
			m.v.SetFlexPanelToAnswers()
		case types.QUIZYESNO:
			// Show Question
			m.v.SetFlexPanelToQuestion(question.YesNoChoiceQuestion.QuestionText)

			//Wait
			time.Sleep(question.TimeToReveal)

			// Update question text
			m.v.UpdateQuestionText(question.YesNoChoiceQuestion.QuestionText)
			m.v.UpdateAnswers(question.TimeToReveal, question.YesNoChoiceQuestion.Answers)
			m.v.SetFlexPanelToAnswers()

		}
	}
	return nil
}
