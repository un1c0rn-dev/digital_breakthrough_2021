package BlackBox

import (
	"fmt"
	"sync"
	"time"
	"unicorn.dev.web-scrap/MagicBox"
)

func Parse(magickRequestPipe chan MagicBox.MagickRequest) {
	for {
		time.Sleep(time.Second)
		fmt.Println("Parsing...")
		var mr MagicBox.MagickRequest
		mr = <-magickRequestPipe
		if len(mr.ContextRequires.Name) != 0 {
			for _, content := range mr.Site.Content {
				pageContext, _ := GetPageContext(content)
				fmt.Println(pageContext)
			}
		}

	}
}

func StartParser(wg *sync.WaitGroup, magickRequestPipe chan MagicBox.MagickRequest) {
	defer wg.Done()
	Parse(magickRequestPipe)
}
