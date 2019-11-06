package main

import (
	BEB "./BEB"
)

type MessageHandlerDaemon struct {
	beb       BEB.BestEffortBroadcast_Module
	host      string
	receivers []string
	_         struct{}
}

func (daemon *MessageHandlerDaemon) SendMessage(msg string) {
	req := BEB.BestEffortBroadcast_Req_Message{
		Addresses: daemon.receivers,
		Message:   msg}
	daemon.beb.Req <- req
}

func (daemon *MessageHandlerDaemon) MessageLoop() {
	/*
		str1 := "MESSAGE_IM_ALIVE;;/data:/;;3"
		res1 := strings.Split(str1, ";;/data:/;;")

		fmt.Println(res1[0])
		fmt.Println(res1[1])
	*/

	for {
		in := <-daemon.beb.Ind
		switch in.Message {
		case MImAlive:
		case MRequestConnectedClients:
		}
	}
}

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
