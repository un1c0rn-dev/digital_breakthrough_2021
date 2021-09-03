package WebApi

type SearchRequest struct {
	Category string `json:"category"`
}

type TaskStatusRequest struct {
	Id uint64 `json:"id"`
}

type ResponseStatus struct {
	Status string `json:"status"`
	ID     uint64 `json:"id"`
}

type ResponseTaskStatus struct {
	Status string `json:"status"`
	ID     uint64 `json:"id"`
}
