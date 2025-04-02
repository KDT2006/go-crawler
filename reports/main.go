package reports

import (
	"fmt"
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
