package MagicBox

import (
	"fmt"
	"sync"
	"time"
	"unicorn.dev.web-scrap/BlackBox"
)

func Scrap(magickRequestPipe chan BlackBox.MagickRequest) {
	for {
		fmt.Println("Scrapping...")
		mr := BlackBox.MagickRequest{
			ContextRequires: BlackBox.ContextRequires{Name: "qwe"},
			Site:            BlackBox.Site{},
		}
		magickRequestPipe <- mr
		time.Sleep(time.Second)
	}
}

func StartScrapper(wg *sync.WaitGroup, magickRequestPipe chan BlackBox.MagickRequest) {
	defer wg.Done()
	Scrap(magickRequestPipe)
}
