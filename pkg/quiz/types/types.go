package types

import (
	"encoding/json"
	"time"
)

// DATA TYPES
const (
	// TYPEUSER defines the data type as user data
	TYPEUSER = "typeUser"

	// TYPEQUIZ defines the data type as quiz data
	TYPEQUIZ = "typeQuiz"

	// TYPESCORE defines the data type as score data
	TYPESCORE = "typeScore"

	// TYPESTATUS defines the data type as quiz status data
	TYPESTATUS = "typeStatus"
)

// QUESTION TYPES
const (
	// QUIZYESNO defines the data type as user data
	QUIZYESNO = "quizYesNo"

	// QUIZMULTI defines the data type as quiz data
	QUIZMULTI = "quizMulti"

	// QUIZRECEIPT means the data is a receipt of a new question
	QUIZRECEIPT = "quizReceipt"

	// QUIZANSWER means the receipt of an answer
	QUIZANSWER = "quizAnswer"
)

// USER TYPES
const (
	// USERANSWER defines the data type as user data
	USERANSWER = "userAnswer"

	// USERJOIN defines the data type as quiz data
	USERJOIN = "userJoin"
)

// DataWrapper is used to wrap data back and forth between client and server and decode to the write type
type DataWrapper struct {
	DataType string          `json:"dataType"`
	Data     json.RawMessage `json:"data"`
}

// ---- Quiz types ---- //

// Quiz defines the questions
type Quiz struct {
	Name      string     `json:"name"`
	Questions []Question `json:"question"`
}

// Question defines the type of question
type Question struct {
	Type                string                  `json:"type"`
	TimeToReveal        time.Duration           `json:"timeToReveal"`
	MultiChoiceQuestion *MultipleChoiceQuestion `json:"multipleChoice,omitempty"`
	YesNoChoiceQuestion *YesNoChoiceQuestion    `json:"yesNoChoice,omitempty"`
}

// MultipleChoiceQuestion defines a question with multiple answers
type MultipleChoiceQuestion struct {
	QuestionText string   `json:"question"`
	Answers      []string `json:"answers"`
	Answer       string   `json:"answer"`
}

// YesNoChoiceQuestion defines a question with multiple answers
type YesNoChoiceQuestion struct {
	QuestionText string   `json:"question"`
	Answers      []string `json:"answers"`
	Answer       string   `json:"answer"`
}

// ---- User types ---- //

// User contains everything needed about a connected user
type User struct {
	Username       string    `json:"username"`
	Address        string    `json:"address"`
	ConnectionTime time.Time `json:"connectionTime"`
	UID            string    `json:"uid"`
}

// ---- Answer types ---- //

// QuizReceipt - sent in receipt of a question
type QuizReceipt struct {
	ReceiptTime time.Time `json:"receiptTime"`
	User        User      `json:"user"`
}

// QuizAnswer - sent when a question is answered
type QuizAnswer struct {
	ReceiptTime time.Time `json:"receiptTime"`
	Answer      string    `json:"answer,omitempty"`
	User        User      `json:"user"`
}

// ---- Status types ---- //

// Status - is used to send status information to a client
type Status struct {
	Error   *string `json:"error:omitempty"`
	Note    *string `json:"error:note"`
	Welcome *string `json:"error:welcome"`
}

// ---- Score types ---- //

// UserScore - is used to manage an indifidual users score
type UserScore struct {
	User  *User
	Score int64
	//FastestAnswer time.Duration
}

//ScoreBoard - is used to manage the scores of the game
type ScoreBoard struct {
	Users         []UserScore
	FastestAnswer int64
	FastestUser   string
}
