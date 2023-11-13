package console

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	pluralDir  = ".plural"
	ConfigName = "console.yml"
)

type VersionedConfig struct {
	ApiVersion string  `yaml:"apiVersion"`
	Kind       string  `yaml:"kind"`
	Spec       *Config `yaml:"spec"`
}

type Config struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

func configFile() string {
	folder, _ := os.UserHomeDir()
	return filepath.Join(folder, pluralDir, ConfigName)
}

func ReadConfig() (conf Config) {
	contents, err := os.ReadFile(configFile())
	if err != nil {
		return
	}

	versioned := &VersionedConfig{Spec: &conf}
	if err = yaml.Unmarshal(contents, versioned); err != nil {
		return Config{}
	}
	return
}
