package WebApi

const TASK_STATUS_DEFERRED string = "deferred"
const TASK_STATUS_ACTIVE string = "active"
const TASK_STATUS_DONE string = "done"
const TASK_STATUS_ERROR string = "error"

type Task struct {
	Id     uint64 `json:"id"`
	Status string `json:"status"`
}

type TaskContext struct {
	IdsCounter uint64
	Tasks      map[uint64]*Task
}

type context struct {
	TaskContext TaskContext
}

var Context context

func initTaskContext() {
	Context.TaskContext.IdsCounter = 0
	Context.TaskContext.Tasks = make(map[uint64]*Task)
}

func InitContext() {
	initTaskContext()
}

func createTaskContext() *Task {
	task := new(Task)
	lastId := &Context.TaskContext.IdsCounter
	*lastId++
	Context.TaskContext.Tasks[*lastId] = task
	task.Id = *lastId
	task.Status = TASK_STATUS_DEFERRED
	return task
}

func getTaskContext(id uint64) *Task {
	if task, ok := Context.TaskContext.Tasks[id]; ok {
		return task
	}

	return nil
}

func removeTaskContext(taskContext *Task) {
	if _, ok := Context.TaskContext.Tasks[taskContext.Id]; ok {
		delete(Context.TaskContext.Tasks, taskContext.Id)
	}
}
