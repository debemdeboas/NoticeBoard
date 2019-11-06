package main

import (
	"fmt"

	BEB "./BEB"
)

type MessageHandlerDaemon struct {
	beb       BEB.BestEffortBroadcast_Module
	host      string
	receivers []string
	_         struct{}
}

func (daemon *MessageHandlerDaemon) SendMessage(msg string) {
	for {
		req := BEB.BestEffortBroadcast_Req_Message{
			Addresses: daemon.receivers,
			Message:   msg}
		daemon.beb.Req <- req
	}
}

func StartMessageHandlerDaemon(addrs []string) MessageHandlerDaemon {
	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}
	beb.Init(addrs[0])
	receivers := addrs[1:]

	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind
			fmt.Printf("Message from %v: %v\n", in.From, in.Message)

		}
	}()

	return MessageHandlerDaemon{}
}
