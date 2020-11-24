package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_MatchService(t *testing.T) {
	svr := Server{Services: []*Service{
		{Name: "foo"}, {Name: "foo.bar.baz"}, {Name: "foo.baz"},
	}}
	s := svr.matchService("foo")
	assert.EqualValues(t, "foo", s.Name)
	s = svr.matchService("foo.bar.baz")
	assert.EqualValues(t, "foo.bar.baz", s.Name)
	s = svr.matchService("foo.baz")
	assert.EqualValues(t, "foo.baz", s.Name)
	s = svr.matchService("foo.bar")
	assert.EqualValues(t, "foo", s.Name)
	s = svr.matchService("foo.baz.bar")
	assert.EqualValues(t, "foo.baz", s.Name)
	s = svr.matchService("bar")
	assert.Nil(t, s)

	svr = Server{Services: []*Service{
		{Name: "foo"}, {Name: "foo.bar.baz"}, {Name: "foo.baz"}, {Name: "*"},
	}}
	s = svr.matchService("bar")
	assert.EqualValues(t, "*", s.Name)
}
