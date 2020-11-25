package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_MatchService(t *testing.T) {
	svr := Server{service: Service{
		"foo": nil, "foo.bar.baz": nil, "foo.baz": nil,
	}}
	name, _ := svr.matchService("foo")
	assert.Equal(t, "foo", name)
	name, _ = svr.matchService("foo.bar.baz")
	assert.Equal(t, "foo.bar.baz", name)
	name, _ = svr.matchService("foo.baz")
	assert.Equal(t, "foo.baz", name)
	name, _ = svr.matchService("foo.bar")
	assert.Equal(t, "foo", name)
	name, _ = svr.matchService("foo.baz.bar")
	assert.Equal(t, "foo.baz", name)
	name, _ = svr.matchService("bar")
	assert.Equal(t, "", name)

	svr = Server{service: Service{
		"foo": nil, "foo.bar.baz": nil, "foo.baz": nil, "*": nil,
	}}
	name, _ = svr.matchService("bar")
	assert.Equal(t, "*", name)
}
