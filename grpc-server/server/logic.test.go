package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGravatar ...
func TestGravatar(t *testing.T) {
	var size int32 = 10
	endpoint := "https://www.gravatar.com/avatar/cf38500a2cd3b6a2c8c1d4d8259e83f8?s=%v"
	email := "kamil@lelonek.me"
	url := gravatar(email, size)
	expected := fmt.Sprintf(endpoint, size)

	assert.Equal(t, url, expected, "URLs are not the same.")
}
