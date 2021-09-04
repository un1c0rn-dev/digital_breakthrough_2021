package WebApi

type SearchRequest struct {
	Keywords []string `json:"keywords"`
}

type TaskStatusRequest struct {
	Id uint64 `json:"id"`
}

type ResponseStatus struct {
	Status string   `json:"status"`
	IDs    []uint64 `json:"ids"`
}

type ResponseTaskStatus struct {
	Status string `json:"status"`
	ID     uint64 `json:"id"`
}
