package Tasks

const TaskStatusDeferred string = "deferred"
const TaskStatusActive string = "active"
const TaskStatusDone string = "done"
const TaskStatusError string = "error"

type Task struct {
	Id       uint64 `json:"id"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
}

type TaskContext struct {
	IdsCounter uint64
	Tasks      map[uint64]*Task
}

type Context struct {
	TaskContext TaskContext
}
