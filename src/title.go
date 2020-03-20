// REFERENCE:
// - https://siongui.github.io/2016/05/10/go-get-html-title-via-net-html/

package main

import (
	"net/http"
	"golang.org/x/net/html"
	"io"
	"os"
	"fmt"
)

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func GetHtmlTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}

func GetTitle(url string)(string, bool) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if title, ok := GetHtmlTitle(resp.Body); ok {
		return title, true
	} else {
		return "", false
	}
}


func GenerateMd(url string, title string) string {
	return "[" + title + "]" + "(" + url + ")"
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		return // no args
	}

	url := args[1]
	title, ok := GetTitle(url)
	if !ok {
		return // error
	}
	md := GenerateMd(url, title)

	fmt.Println(md)
}
