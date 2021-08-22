package MagicBox

import (
	"fmt"
	"sync"
	"time"
)

type ContextRequires struct {
	Name string
}

type Site struct {
	Content []PageContent
}

type MagickRequest struct {
	ContextRequires ContextRequires
	Site            Site
}

func Scrap(magickRequestPipe chan MagickRequest) {
	for {
		fmt.Println("Scrapping...")

		options := SearchOptions{
			SearchLimit:    1,
			TargetLanguage: LANG_CODE_RU,
		}
		f := SearchQuery("картошка", &options)
		results := <-f

		mr := MagickRequest{
			ContextRequires: ContextRequires{
				Name: "кортошка",
			},
			Site: Site{
				Content: []PageContent{},
			},
		}

		if results != nil {
			for i, item := range results {
				pageContent, err := ParsePage(item.URL)
				mr.Site.Content = append(mr.Site.Content, pageContent)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("%d. %s - %s\n", i, item.Title, item.URL)
			}
		}

		go func() { magickRequestPipe <- mr }()
		time.Sleep(time.Second)
	}
}

func StartScrapper(wg *sync.WaitGroup, magickRequestPipe chan MagickRequest) {
	defer wg.Done()
	Scrap(magickRequestPipe)
}
