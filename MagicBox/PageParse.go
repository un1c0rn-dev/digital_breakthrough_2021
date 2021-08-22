package MagicBox

import (
	"container/list"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type WebText struct {
	content string
}

type WebLink struct {
	to      string
	title   string
	content string
}

type PageContent struct {
	text          list.List
	links         list.List
	externalLinks list.List
}

type WebTextError struct {
	s string
}

func (e *WebTextError) Error() string {
	return e.s
}

func isLinkNode(node *html.Node) bool {
	return node.Data == "a" || node.Data == "link"
}

func isScriptNode(node *html.Node) bool {
	return node.Data == "script"
}

func parseTextElemNode(node *html.Node) (WebText, error) {
	webText := WebText{}
	webText.content = strings.Replace(node.Data, "\n", "", -1)
	webText.content = strings.TrimSpace(webText.content)
	if len(webText.content) == 0 {
		return WebText{}, &WebTextError{"Empty tag"}
	}
	//fmt.Println("Parsed text: ", webText)
	return webText, nil
}

func parseLinkElemNode(node *html.Node) WebLink {
	webLink := WebLink{}
	for _, attr := range node.Attr {
		switch attr.Key {
		case "href":
			webLink.to = attr.Val
			break
		case "title":
			webLink.title = attr.Val
			break
		}
	}

	innerNode := node.FirstChild
	if innerNode != nil && innerNode.Type == html.TextNode {
		webLink.content = innerNode.Data
	}

	//fmt.Println("Parsed link: ", webLink)
	return webLink
}

func ParsePage(link string) (PageContent, error) {
	response, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return PageContent{}, err
	}

	if !strings.Contains(response.Status, "200") {
		return PageContent{}, &WebTextError{"Unable to fetch page"}
	}

	htmlPage, err := html.Parse(response.Body)
	if err != nil {
		return PageContent{}, err
	}

	pageContent := PageContent{}

	var processRecursive func(*html.Node)
	processRecursive = func(node *html.Node) {
		switch node.Type {
		case html.TextNode:
			webText, err := parseTextElemNode(node)
			if err != nil {
				pageContent.text.PushBack(webText)
			}
			return
		case html.ElementNode:
			if isScriptNode(node) {
				return
			} else if isLinkNode(node) {
				webLink := parseLinkElemNode(node)
				pageContent.links.PushBack(webLink)
				break
			}
			break
		case html.DocumentNode:
			//println("doc: ", node.Data)
			break
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			processRecursive(c)
		}
	}

	processRecursive(htmlPage)
	return pageContent, nil
}
