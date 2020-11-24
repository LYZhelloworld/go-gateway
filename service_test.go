package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Match(t *testing.T) {
	s := Service{Name: "foo.bar"}

	ok, distance := s.match("foo.bar")
	assert.True(t, ok) // request: foo.bar -> service: foo.bar
	assert.Equal(t, 0, distance)

	ok, distance = s.match("foo.bar.baz")
	assert.True(t, ok) // request: foo.bar.baz -> service: foo.bar
	assert.Equal(t, 1, distance)

	ok, distance = s.match("foo")
	assert.False(t, ok) // request: foo, does not know which "foo.*" should handle

	ok, distance = s.match("foo.baz")
	assert.False(t, ok) // request: foo.baz, wrong name

	s = Service{Name: "*"}

	ok, distance = s.match("foo.bar")
	assert.True(t, ok) // request: foo.bar -> service: *
	assert.Equal(t, 2, distance)

	ok, distance = s.match("foo.bar.baz")
	assert.True(t, ok) // request: foo.bar.baz -> service: *
	assert.Equal(t, 3, distance)

	ok, distance = s.match("foo")
	assert.True(t, ok) // request: foo -> service: *
	assert.Equal(t, 1, distance)
}
