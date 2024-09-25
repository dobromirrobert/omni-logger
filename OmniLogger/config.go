package omnilogger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	LogLevels map[LogLevel]bool `json:"log_levels"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}
	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %v", err)
	}

	return &config, nil
}
