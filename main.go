package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KDT2006/crawler/html"
	"github.com/KDT2006/crawler/reports"
)

type config struct {
	pages              map[string]*reports.Link
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	// Check for proper cli arguments
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		fmt.Println("Usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", os.Args[1])

	baseURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	maxGoRoutines, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	maxPages, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Create initial config
	cfg := &config{
		pages:              make(map[string]*reports.Link),
		baseURL:            baseURL,
		concurrencyControl: make(chan struct{}, maxGoRoutines),
		mu:                 &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
		maxPages:           int(maxPages),
	}

	cfg.wg.Add(1) // Add for the initial goroutine
	go cfg.crawlPage(os.Args[1])

	// Wait for the first goroutine to start
	time.Sleep(time.Millisecond * 10)

	cfg.wg.Wait()

	// Print, Save and Visualize the report
	reports.PrintReports(cfg.pages, cfg.baseURL.String())
	reports.SaveReports(cfg.pages)
	reports.GenerateVisual(cfg.pages, cfg.baseURL.String())
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer func() {
		cfg.wg.Done()            // signal waitgroup
		<-cfg.concurrencyControl // release the spot for another goroutine
	}()

	// Return if length of pages > maxPages
	cfg.mu.Lock()
	if len(cfg.pages) > cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	fmt.Printf("\r\nCrawling page: %s\r\n", rawCurrentURL)

	// Check if rawCurrentURL has the same domain as rawBaseURL check
	baseWithoutPrefix := strings.TrimPrefix(cfg.baseURL.String(), "https://")
	baseDomain := strings.Split(baseWithoutPrefix, "/")[0]

	currentWithoutPrefix := strings.TrimPrefix(rawCurrentURL, "https://")
	currentDomain := strings.Split(currentWithoutPrefix, "/")[0]

	fmt.Printf("Base domain: %s\r\n", baseDomain)
	fmt.Printf("Current domain: %s\r\n", currentDomain)

	// Get the type of link
	var linkType reports.LinkType
	if baseDomain == currentDomain {
		linkType = reports.Internal
	} else {
		linkType = reports.External
	}

	normalizedCurrent, err := html.NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("error normalizeURL(rawCurrentURL) failed:", err)
		return
	}

	// Check if already crawled
	if first := cfg.addPageVisit(normalizedCurrent, linkType); !first {
		fmt.Printf("already crawled %s\r\n", normalizedCurrent)
		return
	}

	// Don't crawl external links
	if baseDomain != currentDomain {
		return
	}

	htmlData, err := html.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("error getHTML(rawCurrentURL) failed:", err)
		return
	}
	fmt.Printf("HTML data for %s:\r\n%s\r\n", rawCurrentURL, htmlData)

	urls, err := html.GetURLsFromHTML(htmlData, cfg.baseURL.String())
	if err != nil {
		fmt.Println("error getURLsFromHTML() failed:", err)
		return
	}

	for _, url := range urls {
		cfg.wg.Add(1) // add to waitgroup
		go func(urlToCrawl string) {
			cfg.concurrencyControl <- struct{}{} // block new goroutines till there's space
			cfg.crawlPage(urlToCrawl)
		}(url)
	}
}

func (cfg *config) addPageVisit(normalizedURL string, linkType reports.LinkType) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, ok := cfg.pages[normalizedURL]; ok {
		cfg.pages[normalizedURL].Count += 1
		return false
	} else {
		cfg.pages[normalizedURL] = reports.NewLink(linkType)
		return true
	}
}
