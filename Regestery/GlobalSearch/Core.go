package main

import (
	"fmt"
	"log"
	"strings"
	"regexp"
	"unicode"
	"net/url"
	"unicorn.dev.web-scrap/MagicBox"

	"github.com/gocolly/colly/v2"
)

const entrypointExp = "картошка оптом купить"
const cacheDir = "./cache"

func main() {
	options := MagicBox.SearchOptions{
		TargetLanguage: MagicBox.LANG_CODE_RU,
		SearchLimit: 10,
	}

	phoneNumRxp := regexp.MustCompile("[+]{0,1}[78][ (-]*[0-9]{3}[ )-]*[0-9]{3}[ -]*[0-9]{2}[ -]*[0-9]{2}")

	log.Print("Asking google for: ", entrypointExp)
	res := <- MagicBox.SearchQuery(entrypointExp, &options)
	for _, element := range res {
		// log.Print("Processing URL: ", element.URL)

		c := colly.NewCollector(
			colly.MaxDepth(1),
			colly.CacheDir(cacheDir),
			// colly.Async(true),
		)
		
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			// log.Print("Downloading: ", link)
			e.Request.Visit(link)
		})

		c.OnXML("//text()", func(e *colly.XMLElement) {
		    text := strings.TrimFunc(e.Text, func(r rune) bool {
		        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		    })

		    text = strings.ReplaceAll(text, "\n", " ")
			text = strings.ReplaceAll(text, "\t", " ")
			text = strings.ReplaceAll(text, "  ", "")

			if(text == "") {
				return 
			}
			if(! strings.Contains(text, "картошка")){
				return
			}

			log.Print(phoneNumRxp.FindAllString(text, 10))

		    // fmt.Println(textFields)
		})

		c.Visit(element.URL)



		fmt.Println(c)

		url.Parse(element.URL)
	}

}