package WebScrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	ipc "unicorn.dev.web-scrap/Ipc"
)

func IpcComm() {
	fmt.Println("Started web scrapper IPC routine.")

	ack := 0
	lastSeq := 0

	for {
		ack++

		var ping = ipc.Ping{
			Ack: ack,
			Seq: lastSeq,
		}

		var serializedPing, _ = json.Marshal(ping)

		var message = ipc.Message{
			To:   "ctx-parser",
			Data: string(serializedPing),
		}

		ipc.Send(message)

		value := ipc.Get(ipc.WebScrapperIpcName)
		if value == nil {
			continue
		}

		stringValue := fmt.Sprintf("%v", value)
		var pong = ipc.Pong{}
		_ = json.Unmarshal([]byte(stringValue), &pong)

		log.Println("Got pong: ", pong)
		lastSeq = pong.Seq

		time.Sleep(time.Second / 2)
	}
}

func Scrap() {
	for {
		time.Sleep(time.Second)
	}
}

func StartScrapper(wg *sync.WaitGroup) {
	defer wg.Done()
	go IpcComm()
	Scrap()
}
