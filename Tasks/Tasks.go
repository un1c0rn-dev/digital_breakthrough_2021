package Tasks

const (
	TaskStatusDeferred string = "deferred"
	TaskStatusActive   string = "active"
	TaskStatusDone     string = "done"
	TaskStatusError    string = "error"
)

const (
	TaskResultReputationUnk  string = "Неизвестно"
	TaskResultReputationBad  string = "Недобросовестная"
	TaskResultReputationMed  string = "Средняя"
	TaskResultReputationGood string = "Хорошая"
)

type TaskResult struct {
	Emails                []string `json:"emails"`
	Phones                []string `json:"phones"`
	ContactPersons        []string `json:"contact_persons"`
	CompanyName           string   `json:"company_name"`
	AverageCapitalization string   `json:"average_capitalization,omitempty"`
	Reputation            string   `json:"reputation,omitempty"`
}

type Task struct {
	Id       uint64       `json:"id"`
	Status   string       `json:"status"`
	Progress string       `json:"progress"`
	Result   []TaskResult `json:"result"`
}

type TaskContext struct {
	IdsCounter uint64
	Tasks      map[uint64]*Task
}

type Context struct {
	TaskContext TaskContext
}
