package BlackBox

import (
	"fmt"
	"sync"
	"time"
)

type ContextRequires struct {
	Name string
}

type Site struct {
}

type MagickRequest struct {
	ContextRequires ContextRequires
	Site            Site
}

func Parse(magickRequestPipe chan MagickRequest) {
	for {
		time.Sleep(time.Second)
		fmt.Println("Parsing...")
		var mr MagickRequest
		mr = <-magickRequestPipe
		if len(mr.ContextRequires.Name) != 0 {
			fmt.Println(mr.ContextRequires.Name)
		}

	}
}

func StartParser(wg *sync.WaitGroup, magickRequestPipe chan MagickRequest) {
	defer wg.Done()
	Parse(magickRequestPipe)
}
