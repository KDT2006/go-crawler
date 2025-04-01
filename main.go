package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Check for proper cli arguments
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", os.Args[1])

	pages := make(map[string]int) // map for storing link counts
	crawlPage(os.Args[1], os.Args[1], pages)

	fmt.Println(pages)
}

func getHTML(raw_url string) (string, error) {
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

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// Check if rawCurrentURL has the same domain as rawBaseURL check
	baseWithoutPrefix := strings.TrimPrefix(rawBaseURL, "https://")
	baseDomain := strings.Split(baseWithoutPrefix, "/")[0]

	currentWithoutPrefix := strings.TrimPrefix(rawCurrentURL, "https://")
	currentDomain := strings.Split(currentWithoutPrefix, "/")[0]

	fmt.Printf("Base domain: %s\r\n", baseDomain)
	fmt.Printf("Current domain: %s\r\n", currentDomain)

	if baseDomain != currentDomain {
		return
	}

	normalizedCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("error normalizeURL(rawCurrentURL) failed:", err)
		return
	}

	// Check if already crawled
	if _, ok := pages[normalizedCurrent]; ok {
		pages[normalizedCurrent] += 1
		return
	} else {
		pages[normalizedCurrent] = 1
	}

	htmlData, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("error getHTML(rawCurrentURL) failed:", err)
		return
	}
	fmt.Printf("HTML data for %s:\r\n%s\r\n", rawCurrentURL, htmlData)

	urls, err := getURLsFromHTML(htmlData, rawBaseURL)
	if err != nil {
		fmt.Println("error getURLsFromHTML() failed:", err)
		return
	}

	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}
