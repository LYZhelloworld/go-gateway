package config

import (
	"testing"

	"github.com/LYZhelloworld/gateway"
	"github.com/stretchr/testify/assert"
)

func TestFromJSON(t *testing.T) {
	cfg := FromJSON([]byte(`{"data":[{"endpoint":"/","method":"GET","service":"test"}]}`))
	assert.Equal(t, "test", cfg[gateway.Endpoint{Path: "/", Method: "GET"}])
}
