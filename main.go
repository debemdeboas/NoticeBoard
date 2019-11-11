package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// IssueID represents an issue's ID
type IssueID struct {
	CreatorID   string
	IssueNumber int
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
	Username   string
	IssueCount int
	Board      *Board
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
		Username:   name,
		IssueCount: 1,
		Board:      &b,
	}
}

func queueUserInput(m MessageHandlerDaemon) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("============ MENU ============\n" +
			"1) Create new issue\n" +
			"2) Edit issue\n" +
			"3) Delete issue\n" +
			"4) Show board...\n" +
			"0) Exit\n> ")
		if scanner.Scan() {
			switch scanner.Text() {
			case "0":
				return
			case "1":
				fmt.Printf("---- Creating a new issue ----\n" +
					"Description: ")
				if scanner.Scan() {
					m.user.CreateIssue(scanner.Text())
					m.SendMessage(MBoard + m.user.Board.toString())
				}
			case "2":
				m.user.EditIssue(
					scanner,
					m.user.Board.GetIssueFromID(
						readIssueID(scanner),
					),
				)
			case "3":
				m.user.DeleteIssue(scanner.Text())
			case "4":
				fmt.Printf("4.1) As text\n4.2) As indented JSON\n> ")
				if scanner.Scan() {
					switch scanner.Text() {
					case "1", "4.1":
						m.user.Board.ShowBoardAsText()
					case "2", "4.2":
						m.user.Board.ShowBoardAsJSON()
					default:
						fmt.Println("Invalid input.")
						log.Info("queueUserInput: ShowBoard: '" + scanner.Text() + "' is not a valid choice.")
					}
				}
			default:
				fmt.Println("Invalid input.")
				log.Info("queueUserInput: '" + scanner.Text() + "' is not a valid choice.")
			}
		}
	}
}

func readIssueID(scanner *bufio.Scanner) (IssueID, error) {
	id := IssueID{}
	fmt.Printf("-------- Select issue --------\n" +
		"Issue creator: ")
	if scanner.Scan() {
		id.CreatorID = scanner.Text()
	}
	fmt.Print("Issue number: ")
	if scanner.Scan() {
		var err error
		id.IssueNumber, err = strconv.Atoi(scanner.Text())
		if err != nil {
			log.Warn("readIssueID: Invalid issue number: " + scanner.Text())
			return id, err
		}
	}
	return id, nil
}

// GetIssueFromID returns the Issue associated with an ID
func (b *Board) GetIssueFromID(id IssueID, err error) *Issue {
	if err != nil {
		return nil
	}

	for i := 0; i < len(b.Issues); i++ {
		if b.Issues[i].ID == id {
			return &(b.Issues[i])
		}
	}

	fmt.Println("Couldn't find issue <" + id.CreatorID + ", " + strconv.Itoa(id.IssueNumber) + ">.")
	return nil
}

// CreateIssue creates an issue and adds that issue to the board
func (u *User) CreateIssue(params string) {
	u.Board.Issues = append(u.Board.Issues,
		Issue{
			ID: IssueID{
				IssueNumber: u.IssueCount,
				CreatorID:   u.Username,
			},
			Content: params,
		},
	)
	u.IssueCount++
	u.Board.TimeStamp = time.Now()
}

// EditIssue edits an issue and saves the changes on the board
func (u *User) EditIssue(scanner *bufio.Scanner, issue *Issue) {
	// For Issue in []Issues -> try to find a specific issue
	u.Board.TimeStamp = time.Now()
}

// DeleteIssue deletes an issue from the board
func (u *User) DeleteIssue(params string) {
	// For Issue in []Issues -> try to find a specific issue
	u.Board.TimeStamp = time.Now()
}

// ShowBoardAsText prints the board to the screen as text
func (b *Board) ShowBoardAsText() {
	for i := 0; i < len(b.Issues); i++ {
		fmt.Printf(b.Issues[i].Content + b.Issues[i].ID.CreatorID + strconv.Itoa(b.Issues[i].ID.IssueNumber))
	}
	fmt.Print("\n")
}

// ShowBoardAsJSON prints the board to the screen as indented JSON
func (b *Board) ShowBoardAsJSON() {
	out, _ := json.MarshalIndent(b, "", "    ")
	fmt.Println(string(out))
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
