package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var configuration *Config = nil

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func GetConfig() Config {
	if configuration == nil {
		conf, err := NewConfig("config.yml")
		if err != nil {
			log.Fatal(err)
		}

		configuration = conf
	}

	return *configuration
}
