package main

import (
	magicBox "unicorn.dev.web-scrap/MagicBox"
)

func main() {

	magicBox.ReqString("qwe",
		[]string{
			"Moscow",
		})
	//startScrapper := flag.Bool("scrapper", false, "Start web scrapper process")
	//startParser := flag.Bool("parser", false, "Start web parser process")

	//flag.Parse()
	//
	//if !*startScrapper && !*startParser {
	//	fmt.Println("Please, choose at least one of mode")
	//	os.Exit(1)
	//}
	//
	//var wg sync.WaitGroup
	//
	//if *startScrapper {
	//	wg.Add(1)
	//	go webScrapper.StartScrapper(&wg)
	//}
	//
	//if *startParser {
	//	wg.Add(1)
	//	go ctxParser.StartParser(&wg)
	//
	//}

	//wg.Wait()
}
