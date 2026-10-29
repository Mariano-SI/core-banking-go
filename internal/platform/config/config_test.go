package config

import (
	"log/slog"
	"os"
	"testing"
)

const testDatabaseURL = "postgres://user:pass@localhost:5432/db"

// unsetEnv removes a variable for the duration of the test.
// t.Setenv registers the cleanup that restores the original value.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
}

func TestLoad_Defaults(t *testing.T) {
	for _, key := range []string{"APP_ENV", "HTTP_PORT", "LOG_LEVEL"} {
		unsetEnv(t, key)
	}
	t.Setenv("DATABASE_URL", testDatabaseURL)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Env != EnvDevelopment {
		t.Errorf("Env = %q, want %q", cfg.Env, EnvDevelopment)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("HTTPPort = %d, want 8080", cfg.HTTPPort)
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Errorf("LogLevel = %v, want %v", cfg.LogLevel, slog.LevelInfo)
	}
	if cfg.HTTPAddr() != ":8080" {
		t.Errorf("HTTPAddr() = %q, want %q", cfg.HTTPAddr(), ":8080")
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("DATABASE_URL", testDatabaseURL)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.IsProduction() {
		t.Errorf("IsProduction() = false, want true")
	}
	if cfg.HTTPPort != 9090 {
		t.Errorf("HTTPPort = %d, want 9090", cfg.HTTPPort)
	}
	if cfg.LogLevel != slog.LevelDebug {
		t.Errorf("LogLevel = %v, want %v", cfg.LogLevel, slog.LevelDebug)
	}
	if cfg.DatabaseURL != testDatabaseURL {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, testDatabaseURL)
	}
}

func TestLoad_InvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{name: "port is not a number", key: "HTTP_PORT", value: "abc"},
		{name: "port out of range", key: "HTTP_PORT", value: "70000"},
		{name: "unknown environment", key: "APP_ENV", value: "staging"},
		{name: "unknown log level", key: "LOG_LEVEL", value: "verbose"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DATABASE_URL", testDatabaseURL)
			t.Setenv(tt.key, tt.value)

			if _, err := Load(); err == nil {
				t.Errorf("Load() with %s=%q: expected an error, got nil", tt.key, tt.value)
			}
		})
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		unsetEnv(t, "DATABASE_URL")
		if _, err := Load(); err == nil {
			t.Error("Load() error = nil, want an error")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "")
		if _, err := Load(); err == nil {
			t.Error("Load() error = nil, want an error")
		}
	})
}
