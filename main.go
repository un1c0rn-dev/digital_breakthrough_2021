package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	ctxParser "unicorn.dev.web-scrap/BlackBox"
	webScrapper "unicorn.dev.web-scrap/MagicBox"
)

func main() {
	startScrapper := flag.Bool("scrapper", false, "Start web scrapper process")
	startParser := flag.Bool("parser", false, "Start web parser process")

	flag.Parse()

	if !*startScrapper && !*startParser {
		fmt.Println("Please, choose at least one of mode")
		os.Exit(1)
	}

	var wg sync.WaitGroup

	if *startScrapper {
		wg.Add(1)
		go webScrapper.StartScrapper(&wg)
	}

	if *startParser {
		wg.Add(1)
		go ctxParser.StartParser(&wg)

	}

	wg.Wait()
}
