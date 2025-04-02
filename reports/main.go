package reports

import (
	"fmt"
	"sort"
)

type kv struct {
	Key   string
	Value int
}

func PrintReports(pages map[string]int, baseURL string) {
	header := fmt.Sprintf(`=============================
  REPORT for %s
=============================`, baseURL)

	fmt.Printf("\n%s\n\n", header)

	sorted := sortReports(pages)
	for _, page := range sorted {
		fmt.Printf("Found %d internal links to %s\n", page.Value, page.Key)
	}
}

func sortReports(pages map[string]int) []kv {
	// Turn the map into a kv slice
	var pagesSlice []kv
	for page, count := range pages {
		pagesSlice = append(pagesSlice, kv{Key: page, Value: count})
	}

	// Sort the slice
	sort.SliceStable(pagesSlice, func(i, j int) bool {
		return pagesSlice[i].Value > pagesSlice[j].Value
	})

	return pagesSlice
}
