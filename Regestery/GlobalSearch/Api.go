package GlobalSearch

import (
	"fmt"
	"log"
	"strings"
	"regexp"
	"unicode"
	"net/url"
	"unicorn.dev.web-scrap/MagicBox"
	"unicorn.dev.web-scrap/Tasks"

	"github.com/gocolly/colly/v2"
)

type Entrypoint struct {

}

type SearchQuery struct {
	Keywords    []string
	Entrypoint	string
	Region      []int
	MaxDomains	[]int
	TimeOutSec 	[]int
}

func NewSearchQuery() SearchQuery {
	return SearchQuery{
		Keywords:   make([]string, 0),
		Entrypoint:	"",
		Region:     make([]int, 0),
		MaxDomains:	make([]int, 0),
		TimeOutSec: make([]int, 0),
	}
}

func Start(query SearchQuery, task *Tasks.Task) {
	tasKStatus := Tasks.TaskStatusDone
	taskProgress := "Готово"
	defer func() {
		task.Status = tasKStatus
		task.Progress = taskProgress
	}()
	task.Status = Tasks.TaskStatusActive

	CreateEntrypointRequest(query)
}

func Stop(task *Tasks.Task) {

}
