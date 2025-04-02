package reports

import (
	"fmt"
	"log"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateBarItems(pages map[string]*Link) []opts.BarData {
	items := make([]opts.BarData, 0)
	for _, link := range pages {
		items = append(items, opts.BarData{
			Value: link.Count,
		})
	}

	return items
}

// GenerateVisual generates a visual plot of all the discovered links and their count
func GenerateVisual(pages map[string]*Link, baseURL string) {
	// Create a new bar instance
	bar := charts.NewBar()

	// Set globals(Title)
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Discovered Links for %s", baseURL),
		}),
		// Don't show the x axis as links are very long to be properly displayed
		charts.WithXAxisOpts(opts.XAxis{
			Show: opts.Bool(false),
		}))

	// Store tooltip data
	keys := []string{}
	// Collect all the pages for the X axis
	for page := range pages {
		keys = append(keys, page)
	}

	// Put data into instance(x axis values-keys are required for tooltips to work!)
	seriesData := generateBarItems(pages)
	bar.SetXAxis(keys).AddSeries("Counts", seriesData)

	// Create a new HTML file for the plot
	f, err := os.Create("visual.html")
	if err != nil {
		log.Println("error creating new file bar.html:", err)
		return
	}
	defer f.Close()

	// Render it!
	if err := bar.Render(f); err != nil {
		log.Println("error rendering the plot:", err)
		return
	}
}
