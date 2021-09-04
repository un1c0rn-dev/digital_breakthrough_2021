package WebApi

import (
	"sync"
	"unicorn.dev.web-scrap/Tasks"
)

var contextMtx sync.RWMutex

var Context Tasks.Context

func initTaskContext() {
	Context.TaskContext.IdsCounter = 0
	Context.TaskContext.Tasks = make(map[uint64]*Tasks.Task)
}

func InitContext() {
	initTaskContext()
}

func createTaskContext() *Tasks.Task {
	task := new(Tasks.Task)
	lastId := &Context.TaskContext.IdsCounter
	*lastId++
	contextMtx.Lock()
	Context.TaskContext.Tasks[*lastId] = task
	contextMtx.Unlock()
	task.Id = *lastId
	task.Status = Tasks.TaskStatusDeferred
	return task
}

func getTaskContext(id uint64) *Tasks.Task {
	contextMtx.RLock()
	if task, ok := Context.TaskContext.Tasks[id]; ok {
		contextMtx.RUnlock()
		return task
	}

	contextMtx.RUnlock()
	return nil
}

func removeTaskContext(taskContext *Tasks.Task) {
	contextMtx.RLock()
	if _, ok := Context.TaskContext.Tasks[taskContext.Id]; ok {
		contextMtx.RUnlock()

		contextMtx.Lock()
		delete(Context.TaskContext.Tasks, taskContext.Id)
		contextMtx.Unlock()
	}
}
