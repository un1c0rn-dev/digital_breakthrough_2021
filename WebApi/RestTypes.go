package WebApi

import "unicorn.dev.web-scrap/Tasks"

type SearchRequest struct {
	Keywords    []string `json:"keywords"`
	FromDateYMD [3]int   `json:"from_date_ymd,omitempty"`
	ToDateYMD   [3]int   `json:"to_date_ymd,omitempty"`
	Region      []int    `json:"region,omitempty"`
	Okpd        string   `json:"okpd,omitempty"`
	Status      int      `json:"status,omitempty"`
	Placing     []int    `json:"placing,omitempty"`
	Etp         []int    `json:"etp,omitempty"`
	MinPrice    int      `json:"min_price,omitempty"`
	MaxPrice    int      `json:"max_price,omitempty"`
	Fz          int      `json:"fz,omitempty"`
	MaxRequests int      `json:"max_requests,omitempty"`
}

type TaskStatusRequest struct {
	Id []uint64 `json:"ids"`
}

type CollectDataRequest struct {
	Ids []uint64 `json:"ids"`
}

type ResponseStatus struct {
	Status string   `json:"status"`
	IDs    []uint64 `json:"ids"`
}

type ResponseTaskStatus struct {
	Status string `json:"status"`
	ID     uint64 `json:"id"`
}

type ResponseCollectData struct {
	Data map[string][]Tasks.TaskResult `json:"data"`
}
