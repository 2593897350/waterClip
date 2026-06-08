package config

import (
	"testing"
)

func TestDefaultConfigValues(t *testing.T) {
	t.Setenv("API_ADDRESS", "")
	t.Setenv("PROCESSOR_BASE_URL", "")

	cfg := Load()

	if cfg.Address != ":8080" {
		t.Fatalf("expected default address :8080, got %s", cfg.Address)
	}
	if cfg.ProcessorBaseURL != "http://localhost:8000" {
		t.Fatalf("expected default processor URL, got %s", cfg.ProcessorBaseURL)
	}
}

func TestConfigReadsEnvironmentOverrides(t *testing.T) {
	t.Setenv("API_ADDRESS", ":9090")
	t.Setenv("PROCESSOR_BASE_URL", "http://processor:9000")

	cfg := Load()

	if cfg.Address != ":9090" {
		t.Fatalf("expected address override, got %s", cfg.Address)
	}
	if cfg.ProcessorBaseURL != "http://processor:9000" {
		t.Fatalf("expected processor URL override, got %s", cfg.ProcessorBaseURL)
	}
}
