package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Issue struct {
	creatorID int
	id        int
}

type Board struct {
	issues *Issue
}

type User struct {
	listen     chan struct{}
	board      Board
	id         int
	issueCount int
}

var users []User

func (i Issue) getIDAsString() string {
	return strconv.Itoa(i.creatorID) + " | " + strconv.Itoa(i.id)
}

func newUser(b Board, id int) User {
	u := User{make(chan struct{}), b, id, 0}

	go func() {
		for {
			<-u.listen
			refresh()
		}
	}()

	return u
}

func refresh() {

}

// global users array

func main() {
	fmt.Print("Number of users: ")
	numUsers, err := bufio.NewReader(os.Stdin).ReadString('\n')

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(numUsers)

	// mainBoard := Board{nil}
	// user1 := newUser(mainBoard, 0)
}
