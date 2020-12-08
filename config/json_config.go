package config

import (
	"encoding/json"

	"github.com/LYZhelloworld/gateway"
)

type jsonConfig struct {
	Data []jsonConfigData `json:"data"`
}

type jsonConfigData struct {
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Service  string `json:"service"`
}

// FromJSON creates gateway.Config from JSON data
func FromJSON(data []byte) gateway.Config {
	jsonCfg := jsonConfig{}
	err := json.Unmarshal(data, &jsonCfg)
	if err != nil {
		panic(err)
	}

	cfg := gateway.Config{}
	for _, d := range jsonCfg.Data {
		cfg.Add(d.Endpoint, d.Method, d.Service)
	}
	return cfg
}
