package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlParser(t *testing.T) {
	_, err := parseURL("https://example.com:9000/org/project")
	assert.NoError(t, err)
}
