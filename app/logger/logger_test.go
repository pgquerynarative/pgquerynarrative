package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewFromConfig_invalidLevel_fallsBackToInfo(t *testing.T) {
	l := NewFromConfig("invalid", false)
	if l.z == nil {
		t.Error("NewFromConfig should produce zerolog backend")
	}
}

func TestNewFromConfig_zerologLogsMessageAndFields(t *testing.T) {
	buf := new(bytes.Buffer)
	l := newFromConfigOutput("info", false, buf)
	l.Info("http request", "method", "GET", "path", "/health", "http.status_code", 200)
	out := buf.String()
	if !strings.Contains(out, "http request") {
		t.Errorf("log should contain message: %s", out)
	}
	if !strings.Contains(out, "GET") || !strings.Contains(out, "/health") || !strings.Contains(out, "200") {
		t.Errorf("log should contain fields: %s", out)
	}
}

func TestNewFromConfig_zerologErrLevel(t *testing.T) {
	buf := new(bytes.Buffer)
	l := newFromConfigOutput("info", false, buf)
	l.Err("error response", "http.status_code", 500, "path", "/api/v1/run")
	out := buf.String()
	if !strings.Contains(out, "error response") || !strings.Contains(out, "500") {
		t.Errorf("Err should log message and fields: %s", out)
	}
}

func TestLogger_Debug_noOpWhenWriterBacked(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(buf)
	l.Debug("should not appear", "x", "y")
	if buf.Len() != 0 {
		t.Errorf("Debug with writer-backed logger should be no-op, got %q", buf.String())
	}
}

func TestLogger_Debug_emitsWhenZerologBacked(t *testing.T) {
	buf := new(bytes.Buffer)
	l := newFromConfigOutput("debug", false, buf)
	l.Debug("debug message", "key", "val")
	out := buf.String()
	if !strings.Contains(out, "debug message") || !strings.Contains(out, "val") {
		t.Errorf("Debug with zerolog should emit: %s", out)
	}
}

func TestLogger_Info_Err_withWriterBacked(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(buf)
	l.Info("hello", "component", "test")
	l.Err("oops", "error", "something failed")
	out := buf.String()
	if !strings.Contains(out, "INF") || !strings.Contains(out, "hello") {
		t.Errorf("Info should write INF and message: %s", out)
	}
	if !strings.Contains(out, "ERR") || !strings.Contains(out, "oops") {
		t.Errorf("Err should write ERR and message: %s", out)
	}
}
