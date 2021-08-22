package BlackBox

import (
	"fmt"
	"github.com/jdkato/prose/v2"
	"unicorn.dev.web-scrap/MagicBox"
)

type PageContext struct {
	tags []string
}

func GetPageContext(content MagicBox.PageContent) (PageContext, error) {
	compiledText := ""
	textElem := content.Text.Front()
	for textElem != nil {
		compiledText += " " + textElem.Value.(MagicBox.WebText).Content
		textElem = textElem.Next()
	}
	doc, err := prose.NewDocument(compiledText, prose.WithExtraction(false))
	if err != nil {
		return PageContext{}, err
	}

	// todo: thief an idea from https://github.com/yash1994/spacy-go (or probably write own solution)
	fmt.Println(doc.Tokens())
	return PageContext{}, nil
}
