package MagicBox

import (
	"container/list"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type WebText struct {
	Content string
}

type WebLink struct {
	To      string
	Title   string
	Content string
}

type PageContent struct {
	Text          list.List
	Links         list.List
	ExternalLinks list.List
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
	webText.Content = strings.Replace(node.Data, "\n", "", -1)
	webText.Content = strings.TrimSpace(webText.Content)
	if len(webText.Content) == 0 {
		return WebText{}, &WebTextError{"Empty tag"}
	}
	//fmt.Println("Parsed Text: ", webText)
	return webText, nil
}

func parseLinkElemNode(node *html.Node) WebLink {
	webLink := WebLink{}
	for _, attr := range node.Attr {
		switch attr.Key {
		case "href":
			webLink.To = attr.Val
			break
		case "Title":
			webLink.Title = attr.Val
			break
		}
	}

	innerNode := node.FirstChild
	if innerNode != nil && innerNode.Type == html.TextNode {
		webLink.Content = innerNode.Data
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
		return PageContent{}, &WebTextError{"Unable To fetch page"}
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
			if err == nil {
				pageContent.Text.PushBack(webText)
			}
			return
		case html.ElementNode:
			if isScriptNode(node) {
				return
			} else if isLinkNode(node) {
				webLink := parseLinkElemNode(node)
				pageContent.Links.PushBack(webLink)
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
