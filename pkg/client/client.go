package client

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thebsdbox/kuiz/pkg/quiz/types"
	"github.com/thebsdbox/kuiz/pkg/ui"
)

// The client package is the "engine" for connecting to the the server
// Handling the incoming data and updating the "view"

var loginAttempt bool

// ClientManager handles all required objects for the client to work
type ClientManager struct {
	v    *ui.ClientView
	c    *websocket.Conn
	u    *types.User
	stop chan bool
	addr string
}

// NewClient will build the client engine
func NewClient(user, address string, v *ui.ClientView) *ClientManager {
	newUser := &types.User{
		Username: user,
	}
	m := ClientManager{
		u:    newUser,
		v:    v,
		stop: make(chan bool, 1),
		addr: address,
	}

	v.AddAnswerHandler(func(answer string) {
		// Package the answer and answer time
		a := types.QuizAnswer{
			Answer:      answer,
			ReceiptTime: time.Now(),
			User:        *m.u,
		}
		dw := types.EncodeData(types.QUIZANSWER, a)
		err := m.c.WriteJSON(dw)
		if err != nil {
			ui.ErrorMessage(err)
		}
		// ->DEBUG
		//		fmt.Println(answer)
	})

	return &m
}

// StartClient will take a new client manager and start connectvitiy
func (m *ClientManager) StartClient() error {

	u := url.URL{Scheme: "ws", Host: m.addr, Path: "/"}
	log.Printf("connecting to %s", u.String())
	var err error
	m.c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		ui.ErrorMessage(err)
		return err
	}
	defer m.c.Close()

	if loginAttempt == false {
		d := types.EncodeData(types.TYPEUSER, m.u)
		m.c.WriteJSON(d)
		loginAttempt = true
	}

	for {
		select {

		case <-m.stop:
			log.Println("interrupt")
			//defer close(m.stopped)
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := m.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return err
			}

		default:
			var data types.DataWrapper
			err := m.c.ReadJSON(&data)
			if err != nil {
				ui.ErrorMessage(err)
				fmt.Println("dafuk")
				m.StopClient()
				return err
			}
			m.processIncoming(&data)

		}

	}
}

// StopClient will close all connections and wait one second before ending
func (m *ClientManager) StopClient() {
	m.v.App.Stop()
	// Wait for all things to be cleared up
	<-time.After(time.Second)
}
