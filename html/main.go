package html

import (
	"fmt"
	"io"
	"net/http"
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

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
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

func GetHTML(raw_url string) (string, error) {
	resp, err := http.Get(raw_url)
	if err != nil {
		return "", err
	}

	// Check for HTTP status error
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error response code: %d", resp.StatusCode)
	}

	// Check for proper content type
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("error invalid content type: %s", resp.Header.Get("Content-Type"))
	}
	defer resp.Body.Close()

	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(htmlData), nil
}
