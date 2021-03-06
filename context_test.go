package gateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockResponseWriter struct {
	mock.Mock
	body       []byte
	statusCode int
}

func (m *mockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *mockResponseWriter) Write(bytes []byte) (int, error) {
	args := m.Called(bytes)
	m.body = bytes
	return args.Int(0), args.Error(1)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
	m.statusCode = statusCode
}

func TestContext_Write(t *testing.T) {
	header := &http.Header{}
	testHeader := &http.Header{}
	testHeader.Add("test", "1")

	rw := &mockResponseWriter{}
	rw.On("Header").Return(*header)
	rw.On("Write", mock.Anything).Return(0, nil)
	rw.On("WriteHeader", mock.Anything).Return(nil)

	c := Context{
		Request:        nil,
		StatusCode:     999,
		Response:       []byte("test"),
		Header:         *testHeader,
		serviceName:    "",
		responseWriter: rw,
		isWritten:      false,
	}
	c.write()

	rw.AssertCalled(t, "Header")
	rw.AssertCalled(t, "Write", []byte("test"))
	rw.AssertCalled(t, "WriteHeader", 999)
	assert.Equal(t, 999, rw.statusCode)
	assert.Equal(t, []byte("test"), rw.body)
	assert.Equal(t, "1", header.Get("Test"))
	assert.True(t, c.isWritten)
}

func TestContext_GetServiceName(t *testing.T) {
	const serviceName = "test.service.name"
	c := Context{serviceName: serviceName}
	assert.EqualValues(t, serviceName, c.GetServiceName())
}
