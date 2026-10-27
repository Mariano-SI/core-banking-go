package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/Mariano-SI/core-banking-go/internal/platform/config"
)

func TestNew_ProductionWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, config.Config{Env: config.EnvProduction, LogLevel: slog.LevelInfo})

	log.Info("hello", "key", "value")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("output is not JSON: %v (%q)", err, buf.String())
	}
	if entry["msg"] != "hello" || entry["key"] != "value" {
		t.Errorf("unexpected entry: %v", entry)
	}
}

func TestNew_DevelopmentWritesText(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, config.Config{Env: config.EnvDevelopment, LogLevel: slog.LevelInfo})

	log.Info("hello")

	if !strings.Contains(buf.String(), "msg=hello") {
		t.Errorf("output = %q, want text format containing msg=hello", buf.String())
	}
}

func TestNew_RespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, config.Config{Env: config.EnvDevelopment, LogLevel: slog.LevelWarn})

	log.Info("should be discarded")

	if buf.Len() != 0 {
		t.Errorf("info message written with level warn: %q", buf.String())
	}
}
