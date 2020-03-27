package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thebsdbox/kuiz/pkg/quiz/types"
)

var clients []*websocket.Conn
var upgrader = websocket.Upgrader{} // use default options

func (m *Manager) server(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	clients = append(clients, c)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	var currentUser *types.User
	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	done := make(chan struct{})
	//defer close(done)

	for {
		// Read from socket
		data := types.DataWrapper{}

		err := c.ReadJSON(&data)
		if err != nil {
			// Check that we've actually registered this user before trying to remove them
			if currentUser != nil {
				m.v.DelUserFromView(fmt.Sprintf("%s/%s", currentUser.Username, ip))

				// TODO - we don't delete an old user

				// err = delUser(currentUser.Username, ip)
				// if err != nil {
				// 	ui.ErrorMessage(err)
				// }
			}
			return
		}
		//fmt.Printf("\n\n\n%s", data)
		// Decode Message
		switch data.DataType {
		case types.TYPEUSER:

			// We should recieve a key at this point
			var newUser types.User
			var dw *types.DataWrapper
			json.Unmarshal(data.Data, &newUser)

			// Process user addition
			currentUser, err = addUser(newUser.Username, ip)
			if err != nil {
				errorString := err.Error()
				status := types.Status{
					Error: &errorString,
				}
				dw = types.EncodeData(types.TYPESTATUS, status)
			} else {
				// Add the user to the view
				m.v.AddUserToView(fmt.Sprintf("%s/%s", currentUser.Username, ip), done)
				dw = types.EncodeData(types.TYPEUSER, currentUser)
			}
			// reply with updated User details
			c.WriteJSON(dw)
			welcome := types.Status{
				Welcome: &m.q.Name,
			}
			dw = types.EncodeData(types.TYPESTATUS, welcome)
			c.WriteJSON(dw)
		case types.QUIZRECEIPT:

			// A User has recieved a question
			var r types.QuizReceipt
			json.Unmarshal(data.Data, &r)
			// Add them to the collection of receipts
			m.qm.receipts = append(m.qm.receipts, r)
		case types.QUIZANSWER:

			// A User has answered a question
			var a types.QuizAnswer
			err = json.Unmarshal(data.Data, &a)
			if err != nil {
				fmt.Println("Error unmarshalling answer")
			}
			m.qm.answers = append(m.qm.answers, a)
		}
	}
	//	}()

	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case <-ticker.C:
	// 		// time out of a second expires
	// 		// wrappedData := types.EncodeData(types.TYPEQUIZ, q.Questions[0])
	// 		// b, _ := json.Marshal(wrappedData)
	// 		// err := c.WriteMessage(websocket.TextMessage, b)
	// 		// if err != nil {
	// 		// 	log.Println("write:", err)
	// 		// 	return
	// 		// }
	// 	}
	// }
}
