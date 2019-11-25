package main

import (
	"strings"
	"time"

	BEB "./BEB"
	log "github.com/sirupsen/logrus"
)

// MessageHandlerDaemon Handler for network messages
type MessageHandlerDaemon struct {
	beb       BEB.BestEffortBroadcast_Module
	host      string
	receivers []string
	user      *User
}

// SendMessage Send messages
func (daemon *MessageHandlerDaemon) SendMessage(msg string) {
	log.Info("Sending message { " + msg + " }")
	req := BEB.BestEffortBroadcast_Req_Message{
		Addresses: daemon.receivers,
		Message:   msg}
	daemon.beb.Req <- req
}

// MessageLoop Handle network messages
func (daemon *MessageHandlerDaemon) MessageLoop() {
	go daemon.SendMessage(MNewUser)

	// Wait for a network message from the other users in order to get the
	// most recent version of the board or create a new local board (in this
	// case this instance is the first one to connect)
	createBoard := false
	select {
	case in := <-daemon.beb.Ind:
		splitMessage := strings.Split(in.Message, MMessageIDDataSeparator)
		messageID := splitMessage[0]
		messageData := splitMessage[1]
		if messageID+MMessageIDDataSeparator == MBoard {
			daemon.HandleBoardF(messageData)
		} else {
			log.Warn("Invalid first message. Got { " + messageID + " } ; Creating new board.")
			createBoard = true
		}
	case <-time.After(CTimeout * time.Millisecond):
		log.Info("Board receival timeout. Creating new board.")
		createBoard = true
	}
	if createBoard {
		daemon.user.Board.TimeStamp = time.Now()
	}

	for {
		in := <-daemon.beb.Ind
		splitMessage := strings.Split(in.Message, MMessageIDDataSeparator)
		messageID := splitMessage[0]
		messageData := splitMessage[1]
		log.Info("Received message:" +
			" { " + messageID + " } " +
			" { " + messageData + " }")
		switch messageID + MMessageIDDataSeparator {
		case MNewUser:
			daemon.HandleNewUser()
		case MBoard:
			daemon.HandleBoard(messageData)
		default:
			log.Warn("MessageHandlerDaemon.MessageLoop: Invalid message: " +
				"ID: { " + messageID + " } " +
				"Data: { " + messageData + " }")
		}
	}
}

// HandleBoardF forcibly updates the current board
func (daemon *MessageHandlerDaemon) HandleBoardF(msg string) {
	daemon.user.Board = buildBoard(msg)
}

// HandleNewUser Handles the New User message (sent when a new user connects to this client)
func (daemon *MessageHandlerDaemon) HandleNewUser() {
	daemon.SendMessage(MBoard + daemon.user.Board.toString())
}

// HandleBoard Updates the user's Board if the received Board's version (timestamp) is higher than
// the user's current Board
func (daemon *MessageHandlerDaemon) HandleBoard(msg string) {
	b := buildBoard(msg)
	if b.TimeStamp.After(daemon.user.Board.TimeStamp) {
		log.Info("HandleBoard: New board received." +
			" Old board: { " + daemon.user.Board.toString() + " } " +
			"New board: { " + msg + " } ")
		daemon.user.Board = b
	}
}

// StartMessageHandlerDaemon Create and start a new Message Handler Daemon
func StartMessageHandlerDaemon(addrs []string, user User) MessageHandlerDaemon {
	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}
	beb.Init(addrs[0])

	daemon := MessageHandlerDaemon{
		beb:       beb,
		host:      addrs[0],
		receivers: addrs[1:],
		user:      &user,
	}

	// Start the daemon's message loop in a separate thread
	go daemon.MessageLoop()

	return daemon
}
