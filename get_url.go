package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func traverseTree(node *html.Node, URLChan chan<- string) {
	// Check if it's an anchor element
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				URLChan <- attr.Val
			}
		}
	}

	// Recursive loop on the child
	if node.FirstChild != nil {
		traverseTree(node.FirstChild, URLChan)
	}

	// Recursive loop on the next sibling
	if node.NextSibling != nil {
		traverseTree(node.NextSibling, URLChan)
	}
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	node, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	// Channel for receiving links
	URLChan := make(chan string)
	var urls []string

	go func() {
		defer close(URLChan)
		traverseTree(node, URLChan)
	}()

	for url := range URLChan {
		// Check for relative paths
		if !strings.Contains(url, "https://") {
			url = fmt.Sprintf("%s%s", rawBaseURL, url)
		}
		urls = append(urls, url)
	}

	return urls, nil
}
