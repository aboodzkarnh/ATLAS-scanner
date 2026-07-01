package core

import (
    "fmt"
    "os"
    "io/ioutil"

    "github.com/goccy/go-yaml"
)

// Config represents the main application configuration
type Config struct {
    Database struct {
        Path string `yaml:"path"`
    } `yaml:"database"`
    Stealth struct {
        RandomAgent bool `yaml:"random_agent"`
        DelayMin    int  `yaml:"delay_min"`
        DelayMax    int  `yaml:"delay_max"`
        Timeout     int  `yaml:"timeout"`
    } `yaml:"stealth"`
    Scan struct {
        Threads  int `yaml:"threads"`
        MaxDepth int `yaml:"max_depth"`
    } `yaml:"scan"`
}

// LoadConfig reads config.yaml from the given path
func LoadConfig(path string) (*Config, error) {
    if path == "" {
        return nil, fmt.Errorf("config path is empty")
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("cannot read config: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }

    // ENV overrides (placeholders - not activated yet)
    if envDB := os.Getenv("ATLAS_DB_PATH"); envDB != "" {
        cfg.Database.Path = envDB
    }

    return &cfg, nil
}
