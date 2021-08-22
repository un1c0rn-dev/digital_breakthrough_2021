package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"unicorn.dev.web-scrap/BlackBox"
	"unicorn.dev.web-scrap/MagicBox"
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
	magickRequestChan := make(chan BlackBox.MagickRequest)

	if *startScrapper {
		wg.Add(1)
		go MagicBox.StartScrapper(&wg, magickRequestChan)
	}

	if *startParser {
		wg.Add(1)
		go BlackBox.StartParser(&wg, magickRequestChan)

	}

	wg.Wait()
}
