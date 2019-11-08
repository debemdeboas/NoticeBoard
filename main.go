package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Issue represents an issue
type Issue struct {
	creatorID int
	id        int
	content   string

	_ struct{}
}

// Board represents an array of issues
type Board struct {
	issues *Issue

	_ struct{}
}

// User represents a user
type User struct {
	board      Board
	id         int
	issueCount int
	username   string

	_ struct{}
}

// Used as a means to pass a board as a message with BEB
func (b Board) toString() string {
	return ""
}

func buildBoard(strBoard string) Board {
	return Board{}
}

func (i Issue) getIDAsString() string {
	return strconv.Itoa(i.creatorID) + " | " + strconv.Itoa(i.id)
}

func newUser(b Board, id int) User {
	u := User{
		board:      b,
		id:         id,
		issueCount: 0,
	}

	return u
}

func queueUserInput() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("1) Create new issue\n2) Edit issue\n3) Delete issue\n0) Exit\n> ")
		if scanner.Scan() {
			switch scanner.Text() {
			case "0":
				return
			case "1":

			case "2":

			case "3":

			default:
				fmt.Println("Invalid input.")
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port")
		return
	}

	addrs := os.Args[1:]
	fmt.Println(addrs)

	fmt.Printf("Hello! How do you wish to be called?\n> ")
	scanner := bufio.NewScanner(os.Stdin)
	var username string
	if scanner.Scan() {
		username = scanner.Text()
	}
	else {
		username = "Guest" + generateRandomID()
	}

	messageDaemon := StartMessageHandlerDaemon(addrs)

	// queueUserInput()
	<-(make(chan struct{}))
}
