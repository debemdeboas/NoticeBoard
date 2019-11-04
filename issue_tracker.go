package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	BEB "./BEB"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addrs := os.Args[1:]
	fmt.Println(addrs)

	receivers := addrs[1:]
	fmt.Println(receivers)

	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	beb.Init(addrs[0])
	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
			}
			req := BEB.BestEffortBroadcast_Req_Message{
				Addresses: receivers,
				Message:   msg}
			beb.Req <- req
		}
	}()

	// receptor de broadcasts
	go func() {
		for {

			in := <-beb.Ind
			fmt.Printf("Message from %v: %v\n", in.From, in.Message)

		}
	}()

	blq := make(chan int)
	<-blq
}
