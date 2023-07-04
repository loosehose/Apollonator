package pkg

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKey       string `yaml:"api_key"`
	Organization string `yaml:"organization"`
	Email        bool   `yaml:"email"`
	Title        bool   `yaml:"title"`
}

// ParseYAML reads and parses a YAML configuration file
func ParseYAML(configFile string) (Config, error) {
	config := struct {
		Apollonator Config `yaml:"apollonator"`
	}{}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, err
	}

	// Validate Configurations
	if config.Apollonator.APIKey == "" {
		return Config{}, errors.New("missing API key in the configuration")
	}

	if config.Apollonator.Organization == "" {
		return Config{}, errors.New("missing Organization in the configuration")
	}

	return config.Apollonator, nil
}
