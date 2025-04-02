# Go-Crawler

This project is a **concurrent web crawler** written in Go. It was designed to explore and extract links from web pages, classify them as internal or external, and generate reports and visualizations. The project was both fun and educational, providing hands-on experience with Go's concurrency model, HTML parsing, and data visualization.

## Features

- **Concurrent Crawling**: Efficiently crawls web pages using goroutines and channels.
- **Link Classification**: Differentiates between internal and external links.
- **HTML Parsing**: Extracts links from HTML content.
- **Reports**:
    - Generates a CSV report of discovered links.
    - Prints a summary report to the console.
    - Creates a visual bar chart of link counts using `go-echarts`.

## Usage

### Prerequisites

- Go 1.24 or later installed on your system.

### Running the Crawler

```bash
go run main.go <URL> <maxConcurrency> <maxPages>
```

- `<URL>`: The starting URL for the crawl.
- `<maxConcurrency>`: Maximum number of concurrent goroutines.
- `<maxPages>`: Maximum number of pages to crawl.

### Example

```bash
go run main.go https://example.com 10 100
```

### Output

- **Console Report**: Displays a summary of discovered links.
- **CSV Report**: Saved as `report.csv`.
- **Visualization**: Saved as `visual.html`.

## Testing

Unit tests are provided for the `html` package. To run the tests:

```bash
go test ./...
```

## Dependencies

- [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html): HTML parsing.
- [github.com/stretchr/testify](https://pkg.go.dev/github.com/stretchr/testify): Testing utilities.
- [github.com/go-echarts/go-echarts/v2](https://pkg.go.dev/github.com/go-echarts/go-echarts/v2): Data visualization.

## License

This project is open-source and available under the [MIT License](LICENSE).

---
Happy crawling! 🚀