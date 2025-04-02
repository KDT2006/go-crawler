package html

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURLsFromHTML(t *testing.T) {
	// Test: Simple test
	inputURL := "https://blog.boot.dev"
	inputBody := `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`
	expected := []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"}
	actual, err := GetURLsFromHTML(inputBody, inputURL)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expected, actual))

	// Test: Multiple absoulte paths
	inputURL = "https://blog.boot.dev"
	inputBody = `
<html>
	<body>
		<a href="https://example.com/article/1">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://example.com/auth/user">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`
	expected = []string{
		"https://example.com/article/1",
		"https://other.com/path/one",
		"https://example.com/auth/user",
	}
	actual, err = GetURLsFromHTML(inputBody, inputURL)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expected, actual))

	// Test: Empty href
	inputURL = "https://blog.boot.dev"
	inputBody = `
<html>
	<body>
		<a>
			<span>Boot.dev</span>
		</a>
		<a>
			<span>Boot.dev</span>
		</a>
		<a>
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`
	actual, err = GetURLsFromHTML(inputBody, inputURL)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(actual))

	// Test: Empty body
	inputURL = "https://blog.boot.dev"
	inputBody = ``
	expected = []string{}
	actual, err = GetURLsFromHTML(inputBody, inputURL)
	assert.Nil(t, err)
}
