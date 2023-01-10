package config

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Sensor is the output structure for ESP8266
type Sensor struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	AirQuality  string `json:"air"`
	Light       string `json:"light"`
}

// Config is structure of config file
type Config struct {
	Clients []Client `yaml:"clients"`
}

// Client is array of clients in config file
type Client struct {
	Name string `yaml:"name"`
	IP   string `yaml:"ip"`
}

// Load reads YAML from reader and unmarshal in Config
func Load(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
