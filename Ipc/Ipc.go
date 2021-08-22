package Ipc

import (
	"container/list"
	"sync"
)

type Ping struct {
	Ack int `json:"ack"`
	Seq int `json:"seq"`
}

type Pong struct {
	Ack int `json:"ack"`
	Seq int `json:"seq"`
}

type Message struct {
	To   string
	Data string
}

var _ipcQueueMap map[string]*list.List
var _queueMutex sync.Mutex

func init() {
	_ipcQueueMap = make(map[string]*list.List)
}

func Send(message Message) {
	_queueMutex.Lock()
	if _, ok := _ipcQueueMap[message.To]; !ok {
		_ipcQueueMap[message.To] = new(list.List)
	}
	_ipcQueueMap[message.To].PushBack(message.Data)
	_queueMutex.Unlock()
}

func Get(receiver string) interface{} {
	_queueMutex.Lock()
	if queue, ok := _ipcQueueMap[receiver]; ok {
		if queue.Len() == 0 {
			_queueMutex.Unlock()
			return nil
		} else {
			valuePrt := queue.Front()
			value := valuePrt.Value
			queue.Remove(valuePrt)
			_queueMutex.Unlock()
			return value
		}
	}
	_queueMutex.Unlock()
	return nil
}
