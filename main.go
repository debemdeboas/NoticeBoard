package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// IssueID represents an issue's ID
type IssueID struct {
	IssueNumber int
	CreatorID   string
}

// Issue represents an issue
type Issue struct {
	ID      IssueID
	Content string
}

// Board represents an array of issues
type Board struct {
	Issues    []Issue
	TimeStamp time.Time
}

// User represents a user
type User struct {
	Board      *Board
	IssueCount int
	Username   string
}

func buildBoard(jBoard string) *Board {
	var b *Board
	if err := json.Unmarshal([]byte(jBoard), &b); err != nil {
		log.Fatal(err)
	}
	return b
}

func (b *Board) toString() string {
	var buf = new(bytes.Buffer)
	json.NewEncoder(buf).Encode(b)
	return buf.String()
}

// Join combines two boards
func (b *Board) Join(newBoard Board) {
	b.Issues = append(b.Issues, newBoard.Issues...)

	b.TimeStamp = time.Now()
}

func newUser(b Board, name string) User {
	return User{
		Board:      &b,
		Username:   name,
		IssueCount: 0,
	}
}

func queueUserInput(m MessageHandlerDaemon) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("1) Create new issue\n2) Edit issue\n3) Delete issue\n0) Exit\n> ")
		if scanner.Scan() {
			switch scanner.Text() {
			case "0":
				return
			case "1":
				fmt.Printf("--- Creating a new issue ---\n" +
					"Description: ")
				if scanner.Scan() {
					m.user.CreateIssue(scanner.Text())
					m.SendMessage(MBoard + m.user.Board.toString())
				}
			case "2":
				m.user.EditIssue(scanner.Text())
			case "3":
				m.user.DeleteIssue(scanner.Text())
			default:
				fmt.Println("Invalid input.")
				log.Info("queueUserInput: '" + scanner.Text() + "' is not a valid choice.")
			}
		}
	}
}

// CreateIssue creates an issue and adds that issue to the board
func (u *User) CreateIssue(params string) {
	u.Board.Issues = append(u.Board.Issues,
		Issue{
			ID: IssueID{
				u.IssueCount,
				u.Username,
			},
			Content: params,
		},
	)
	u.IssueCount++
	u.Board.TimeStamp = time.Now()
}

// EditIssue edits an issue and saves the changes on the board
func (u *User) EditIssue(params string) {
	// For Issue in []Issues -> try to find a specific issue
	u.Board.TimeStamp = time.Now()
}

// DeleteIssue deletes an issue from the board
func (u *User) DeleteIssue(params string) {
	// For Issue in []Issues -> try to find a specific issue
	u.Board.TimeStamp = time.Now()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port")
		return
	}

	// logFile := time.Now().Format("20060102") + "_IssueTracker.log"
	// file, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	fmt.Println("Could not create logs file. Error: " + err.Error())
	// } else {
	// 	defer file.Close()
	// 	log.SetOutput(file)
	// }

	addrs := os.Args[1:]
	log.Info(addrs)

	fmt.Printf("Hi! How do you wish to be called?\n> ")
	scanner := bufio.NewScanner(os.Stdin)

	// username := "Guest " + strconv.Itoa(rand.Intn(1000))
	var username string
	if scanner.Scan() {
		username = scanner.Text()
	}
	fmt.Println("Hello, " + username + ", and welcome to the Notice Board.\n" +
		"Here you can create, edit and delete issues and synchronize your Board with other users automatically!")

	messageDaemon := StartMessageHandlerDaemon(
		addrs,
		newUser(
			Board{Issues: make([]Issue, 0)},
			username,
		),
	)

	// var buf = new(bytes.Buffer)
	// enc := json.NewEncoder(buf)
	// enc.Encode(messageDaemon.user)
	// f, err := os.Create("daemon.json")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer f.Close()
	// io.Copy(f, buf)

	queueUserInput(messageDaemon)
}
