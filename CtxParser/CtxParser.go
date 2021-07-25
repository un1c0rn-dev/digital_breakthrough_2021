package CtxParser

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	ipc "unicorn.dev.web-scrap/Ipc"
)

func IpcComm() {

	fmt.Println("Started ctx scanner IPC routine.")

	seq := 0
	lastAck := 0

	for {
		value := ipc.Get(ipc.CtxParserIpcName)
		if value == nil {
			continue
		}

		stringValue := fmt.Sprintf("%v", value)

		var ping = ipc.Ping{}
		_ = json.Unmarshal([]byte(stringValue), &ping)

		log.Println("Got ping: ", ping)

		var pong = ipc.Pong{
			Ack: lastAck,
			Seq: seq,
		}

		serializedPong, _ := json.Marshal(pong)

		message := ipc.Message{
			To:   ipc.WebScrapperIpcName,
			Data: string(serializedPong),
		}

		ipc.Send(message)
		seq++
		lastAck = ping.Ack

		time.Sleep(time.Second / 2)
	}

}

func Parse() {
	for {
		time.Sleep(time.Second)
	}
}

func StartParser(wg *sync.WaitGroup) {
	defer wg.Done()
	go IpcComm()
	Parse()
}
