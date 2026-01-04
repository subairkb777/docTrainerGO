package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	InputType string `yaml:"input_type"`
	PDF       struct {
		Path          string `yaml:"path"`
		ExtractImages bool   `yaml:"extract_images"`
	} `yaml:"pdf"`
	Markdown struct {
		Directory    string   `yaml:"directory"`
		AutoDiscover bool     `yaml:"auto_discover"`
		Files        []string `yaml:"files"`
	} `yaml:"markdown"`
	Output struct {
		Directory string `yaml:"directory"`
		Title     string `yaml:"title"`
	} `yaml:"output"`
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Ollama struct {
		Enabled bool   `yaml:"enabled"`
		URL     string `yaml:"url"`
		Model   string `yaml:"model"`
	} `yaml:"ollama"`
}

// Load reads and parses the configuration file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	if config.Output.Directory == "" {
		config.Output.Directory = "docs"
	}
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Ollama.URL == "" {
		config.Ollama.URL = "http://localhost:11434"
	}
	if config.Ollama.Model == "" {
		config.Ollama.Model = "llama3.2"
	}

	return &config, nil
}
