package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeURL(t *testing.T) {
	// Test: remove scheme
	input := "https://blog.boot.dev/path"
	expected := "blog.boot.dev/path"
	actual, err := normalizeURL(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	input = "http://blog.boot.dev/path"
	expected = "blog.boot.dev/path"
	actual, err = normalizeURL(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	// Test: remove trailing slash
	input = "blog.boot.dev/path/"
	expected = "blog.boot.dev/path"
	actual, err = normalizeURL(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
