package config

import (
	"strings"
	"testing"
)

func TestLoadReturnsConfig(t *testing.T) {
	t.Setenv("API_KEY", " test-key ")
	t.Setenv("BASE_URL", " https://example.com/v1 ")
	t.Setenv("MODEL_NAME", " test-model ")

	conf, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if conf.APIKey != "test-key" {
		t.Fatalf("APIKey = %q, want %q", conf.APIKey, "test-key")
	}
	if conf.BaseURL != "https://example.com/v1" {
		t.Fatalf("BaseURL = %q, want %q", conf.BaseURL, "https://example.com/v1")
	}
	if conf.ModelName != "test-model" {
		t.Fatalf("ModelName = %q, want %q", conf.ModelName, "test-model")
	}
}

func TestLoadReportsMissingRequiredConfig(t *testing.T) {
	t.Setenv("API_KEY", " ")
	t.Setenv("BASE_URL", "")
	t.Setenv("MODEL_NAME", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want missing config error")
	}

	msg := err.Error()
	if !strings.Contains(msg, "API_KEY") {
		t.Fatalf("error %q does not mention API_KEY", msg)
	}
	if !strings.Contains(msg, "MODEL_NAME") {
		t.Fatalf("error %q does not mention MODEL_NAME", msg)
	}
}
