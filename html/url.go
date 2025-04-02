package html

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeURL(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", parsed.Host, strings.TrimSuffix(parsed.Path, "/")), nil
}
