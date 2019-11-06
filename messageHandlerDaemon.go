package main

import (
	"strings"

	BEB "./BEB"
)

// MessageHandlerDaemon Handler for network messages
type MessageHandlerDaemon struct {
	beb       BEB.BestEffortBroadcast_Module
	host      string
	receivers []string
	_         struct{}
}

// SendMessage Send messages
func (daemon *MessageHandlerDaemon) SendMessage(msg string) {
	req := BEB.BestEffortBroadcast_Req_Message{
		Addresses: daemon.receivers,
		Message:   msg}
	daemon.beb.Req <- req
}

// MessageLoop Handle network messages
func (daemon *MessageHandlerDaemon) MessageLoop() {
	for {
		in := <-daemon.beb.Ind
		splitMessage := strings.Split(in.Message, MMessageIDAndDataSeparator)
		messageID := splitMessage[0]
		// messageData := splitMessage[1]

		switch messageID {
		case MImAlive:
		case MRequestConnectedClients:
		}
	}
}

// StartMessageHandlerDaemon Create and start a new Message Handler Daemon
func StartMessageHandlerDaemon(addrs []string) MessageHandlerDaemon {
	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}
	beb.Init(addrs[0])

	daemon := MessageHandlerDaemon{
		beb:       beb,
		host:      addrs[0],
		receivers: addrs[1:]}

	return daemon
}
