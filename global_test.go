package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPath(t *testing.T) {
	assert.True(t, isValidPath("/foo/bar/baz"))
	assert.True(t, isValidPath("/foo"))
	assert.False(t, isValidPath("/foo/"))
	assert.True(t, isValidPath("/"))
	assert.False(t, isValidPath("/*")) // not a valid path, but a valid prefix
	assert.True(t, isValidPath(trimPrefix("/*")))
	assert.False(t, isValidPath("/foo*"))
	assert.False(t, isValidPath("/foo/*")) // not a valid path, but a valid prefix
	assert.True(t, isValidPath(trimPrefix("/foo/*")))
}

func TestIsValidService(t *testing.T) {
	assert.True(t, isValidService("foo.bar.baz"))
	assert.True(t, isValidService("foo"))
	assert.False(t, isValidService("foo."))
	assert.False(t, isValidService(""))
	assert.True(t, isValidService("*"))
	assert.False(t, isValidService("foo*"))
	assert.False(t, isValidService(".foo"))
	assert.False(t, isValidService(".foo."))
}

func TestRemoveLastDir(t *testing.T) {
	assert.Equal(t, "/foo/bar", removeLastDir("/foo/bar/baz"))
	assert.Equal(t, "/", removeLastDir("/foo"))
	assert.Equal(t, "", removeLastDir("/"))
}

func TestRemoveLastSubService(t *testing.T) {
	assert.Equal(t, "foo.bar", removeLastSubService("foo.bar.baz"))
	assert.Equal(t, "", removeLastSubService("foo"))
}
