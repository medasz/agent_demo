package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APIKey    string
	BaseURL   string
	ModelName string
}

func Load() (Config, error) {
	conf := Config{
		APIKey:    strings.TrimSpace(os.Getenv("API_KEY")),
		BaseURL:   strings.TrimSpace(os.Getenv("BASE_URL")),
		ModelName: strings.TrimSpace(os.Getenv("MODEL_NAME")),
	}

	var missing []string
	if conf.APIKey == "" {
		missing = append(missing, "API_KEY")
	}
	if conf.ModelName == "" {
		missing = append(missing, "MODEL_NAME")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}

	return conf, nil
}
