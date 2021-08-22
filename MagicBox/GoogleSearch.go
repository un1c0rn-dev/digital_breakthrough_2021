package MagicBox

import (
	"context"
	"fmt"
	googlesearch "github.com/rocketlaunchr/google-search"
)

const LANG_CODE_RU string = "ru"
const LANG_CODE_EN string = "en"

type SearchOptions struct {
	TargetLanguage string
	SearchLimit    uint
}

func SearchQuery(queryString string, options *SearchOptions) chan []googlesearch.Result {

	var ctx = context.Background()
	opts := googlesearch.SearchOptions{
		Limit: 20,
	}

	if options != nil {
		if len(options.TargetLanguage) > 0 {
			opts.LanguageCode = options.TargetLanguage
		}

		if options.SearchLimit > 0 {
			opts.Limit = int(options.SearchLimit)
		}
	}

	f := make(chan []googlesearch.Result)

	go func() {
		results, err := googlesearch.Search(ctx, queryString, opts)
		if err != nil {
			fmt.Errorf("Something went wrong: %v", err)
			f <- nil
		}

		if len(results) == 0 {
			fmt.Errorf("No results returned: %v", results)
			f <- nil
		}

		f <- results
	}()

	return f
}
