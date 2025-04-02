package reports

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
)

// LinkType "enum"
type LinkType int

const (
	Internal = iota
	External
)

// Link holds each link count along with its type
type Link struct {
	LinkType LinkType
	Count    int
}

func NewLink(linkType LinkType) *Link {
	return &Link{
		LinkType: linkType,
		Count:    1,
	}
}

type kv struct {
	Key   string
	Value *Link
}

func PrintReports(pages map[string]*Link, baseURL string) {
	header := fmt.Sprintf(`=============================
  REPORT for %s
=============================`, baseURL)

	fmt.Printf("\n%s\n\n", header)

	sorted := sortReports(pages)
	for _, page := range sorted {
		var linkTypeText string
		if page.Value.LinkType == 0 {
			linkTypeText = "internal"
		} else {
			linkTypeText = "external"
		}
		fmt.Printf("Found %d %s links to %s\n", page.Value.Count, linkTypeText, page.Key)
	}
}

func SaveReports(pages map[string]*Link) {
	file, err := os.Create("report.csv")
	if err != nil {
		log.Println("error creating file report.csv:", err)
		return
	}
	defer file.Close()

	records := [][]string{}

	// Fill in the headers
	records = append(records, []string{"link_type", "link", "count"})

	// Populate the records to write: linkType, link, count
	for page, link := range pages {
		var linkTypeText string
		if link.LinkType == 0 {
			linkTypeText = "internal"
		} else {
			linkTypeText = "external"
		}

		temp := []string{}
		temp = append(temp, linkTypeText)
		temp = append(temp, page)
		temp = append(temp, fmt.Sprintf("%d", link.Count))

		records = append(records, temp)
	}

	// Print to stdout
	w := csv.NewWriter(file)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Flush out any remaining buffer
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Report saved to report.csv!")
}

func sortReports(pages map[string]*Link) []kv {
	// Turn the map into a kv slice
	var pagesSlice []kv
	for page, link := range pages {
		pagesSlice = append(pagesSlice, kv{Key: page, Value: link})
	}

	// Sort the slice
	sort.SliceStable(pagesSlice, func(i, j int) bool {
		return pagesSlice[i].Value.Count > pagesSlice[j].Value.Count
	})

	return pagesSlice
}
