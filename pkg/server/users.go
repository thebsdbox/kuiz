package server

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/thebsdbox/kuiz/pkg/quiz/types"
)

// Users holds all users of the quiz
var Users []types.User

func addUser(username, address string) (*types.User, error) {
	_, found := findUser(username)
	// If the username already exists in the quiz return an error
	if found == true {
		return nil, fmt.Errorf("Username [ %s ] has already registered for this quiz", username)
	}

	// ->DEBUG
	//	fmt.Println("USER GEN")

	// Generate a uid for this user
	UID := uuid.NewV4()
	// if err != nil {
	// 	return nil, err
	// }

	generatedUID := uuid.NewV5(UID, username)

	newUser := types.User{
		Username:       username,
		UID:            generatedUID.String(),
		Address:        address,
		ConnectionTime: time.Now(),
	}

	Users = append(Users, newUser)
	return &newUser, nil
}

func delUser(username, address string) error {
	i, found := findUser(username)
	// If the username already exists in the quiz return an error
	if found == true {
		Users = append(Users[:i], Users[i+1:]...)
		return nil
	}
	return fmt.Errorf("Unable to find user [%s]", username)
}

func generateScores() (scoretext string) {
	return ""
}

func findUser(username string) (index int, found bool) {
	var user types.User
	// Look through the Users and find if the users exists, and where they exist
	for index, user = range Users {
		if username == user.Username {
			found = true
			return
		}
	}
	return
}
