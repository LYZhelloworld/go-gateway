package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterConfig_SetService(t *testing.T) {
	testService := &serviceInfo{}
	rc := routerConfig{}

	ok := rc.setService(http.MethodGet, testService)
	assert.True(t, ok)
	assert.Equal(t, testService, rc.get)
	rc.get = nil

	ok = rc.setService(http.MethodPost, testService)
	assert.True(t, ok)
	assert.Equal(t, testService, rc.post)
	rc.post = nil

	ok = rc.setService(http.MethodPut, testService)
	assert.True(t, ok)
	assert.Equal(t, testService, rc.put)
	rc.put = nil

	ok = rc.setService(http.MethodDelete, testService)
	assert.True(t, ok)
	assert.Equal(t, testService, rc.delete)
	rc.delete = nil

	ok = rc.setService(http.MethodConnect, testService)
	assert.False(t, ok)
}

func TestRouterConfig_GetService(t *testing.T) {
	testService := &serviceInfo{}
	rc := routerConfig{}

	rc.get = testService
	s, ok := rc.getService(http.MethodGet)
	assert.True(t, ok)
	assert.Equal(t, testService, s)
	rc.get = nil

	rc.post = testService
	s, ok = rc.getService(http.MethodPost)
	assert.True(t, ok)
	assert.Equal(t, testService, s)
	rc.post = nil

	rc.put = testService
	s, ok = rc.getService(http.MethodPut)
	assert.True(t, ok)
	assert.Equal(t, testService, s)
	rc.put = nil

	rc.delete = testService
	s, ok = rc.getService(http.MethodDelete)
	assert.True(t, ok)
	assert.Equal(t, testService, s)
	rc.delete = nil

	_, ok = rc.getService(http.MethodConnect)
	assert.False(t, ok)
}
