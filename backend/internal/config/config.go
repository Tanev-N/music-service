package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	Storage StorageConfig `yaml:"storage"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Port string `yaml:"port"`
}

type StorageConfig struct {
	TracksDir     string   `yaml:"tracks_dir"`
	MaxFileSizeMB int      `yaml:"max_file_size_mb"`
	AllowedTypes  []string `yaml:"allowed_types"`
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
