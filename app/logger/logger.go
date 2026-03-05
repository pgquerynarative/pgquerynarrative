// Package logger provides structured console logging with timestamp, level (INF/WRN/ERR),
// message, and optional key=value pairs so all errors and events are visible in a consistent format.
// When created with NewFromConfig (LOG_LEVEL, LOG_PRETTY), uses zerolog for leveled, optional
// colorful output; otherwise uses the built-in format.
package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Level is the log level.
type Level string

const (
	LevelInfo Level = "INF"
	LevelWarn Level = "WRN"
	LevelErr  Level = "ERR"
)

// Logger writes structured log lines: "3:04AM LVL message key1=value1 key2=value2",
// or when zerolog-backed (NewFromConfig), uses zerolog format with optional pretty/color output.
type Logger struct {
	w io.Writer
	z *zerolog.Logger // when set, use zerolog for all output
}

// Default returns a logger that writes to os.Stdout with the built-in format.
func Default() *Logger {
	return &Logger{w: os.Stdout}
}

// New returns a logger that writes to w with the built-in format.
func New(w io.Writer) *Logger {
	return &Logger{w: w}
}

// NewFromConfig builds a zerolog-backed logger from level (e.g. "info", "debug") and pretty (colorful console).
// Invalid level falls back to info. Use for application logging when LOG_LEVEL and LOG_PRETTY are set.
func NewFromConfig(level string, pretty bool) *Logger {
	return newFromConfigOutput(level, pretty, os.Stdout)
}

// newFromConfigOutput builds a zerolog-backed logger writing to w. Used by NewFromConfig and tests.
func newFromConfigOutput(level string, pretty bool, w io.Writer) *Logger {
	lvl := zerolog.InfoLevel
	if level != "" {
		if p, err := zerolog.ParseLevel(level); err == nil {
			lvl = p
		}
	}
	out := io.Writer(w)
	if pretty {
		out = zerolog.ConsoleWriter{Out: w}
	}
	z := zerolog.New(out).With().Timestamp().Logger().Level(lvl)
	return &Logger{z: &z}
}

// Log writes a line with the given level, message, and key-value pairs (alternating key, value).
// Keys are quoted if they contain spaces; values are space-escaped for readability.
func (l *Logger) Log(level Level, msg string, kv ...interface{}) {
	if l.z != nil {
		l.logZerolog(level, msg, kv)
		return
	}
	line := l.format(level, msg, kv)
	_, _ = fmt.Fprintln(l.w, line)
}

func (l *Logger) logZerolog(level Level, msg string, kv []interface{}) {
	var evt *zerolog.Event
	switch level {
	case LevelErr:
		evt = l.z.Error()
	case LevelWarn:
		evt = l.z.Warn()
	default:
		evt = l.z.Info()
	}
	for i := 0; i+1 < len(kv); i += 2 {
		k := fmt.Sprint(kv[i])
		v := kv[i+1]
		if err, ok := v.(error); ok && k == "error" {
			evt = evt.Err(err)
		} else {
			evt = evt.Interface(k, v)
		}
	}
	evt.Msg(msg)
}

func (l *Logger) format(level Level, msg string, kv []interface{}) string {
	ts := time.Now().Format("3:04AM")
	parts := []string{ts, string(level), msg}
	for i := 0; i+1 < len(kv); i += 2 {
		k := fmt.Sprint(kv[i])
		v := fmt.Sprint(kv[i+1])
		if strings.ContainsAny(k, " \t=") {
			k = fmt.Sprintf("%q", k)
		}
		if strings.ContainsAny(v, " \t\n") {
			v = strings.ReplaceAll(v, "\n", " ")
			v = strings.TrimSpace(v)
			if len(v) > 128 {
				v = v[:128] + "..."
			}
			v = fmt.Sprintf("%q", v)
		}
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, " ")
}

// Debug logs at debug level (zerolog only; no-op when using built-in format).
func (l *Logger) Debug(msg string, kv ...interface{}) {
	if l.z != nil {
		evt := l.z.Debug()
		for i := 0; i+1 < len(kv); i += 2 {
			k := fmt.Sprint(kv[i])
			v := kv[i+1]
			if err, ok := v.(error); ok && k == "error" {
				evt = evt.Err(err)
			} else {
				evt = evt.Interface(k, v)
			}
		}
		evt.Msg(msg)
	}
}

// Info logs at INF level.
func (l *Logger) Info(msg string, kv ...interface{}) {
	l.Log(LevelInfo, msg, kv...)
}

// Infof logs at INF level with a format string.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn logs at WRN level.
func (l *Logger) Warn(msg string, kv ...interface{}) {
	l.Log(LevelWarn, msg, kv...)
}

// Warnf logs at WRN level with a format string.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Err logs at ERR level.
func (l *Logger) Err(msg string, kv ...interface{}) {
	l.Log(LevelErr, msg, kv...)
}

// Errf logs at ERR level with a format string.
func (l *Logger) Errf(format string, args ...interface{}) {
	l.Err(fmt.Sprintf(format, args...))
}

// defaultLogger is the package-level logger used by apilog and optional callers.
var defaultLogger = Default()

// SetDefault sets the logger used by package-level apilog. Main can call this to use a test buffer.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// DefaultLogger returns the current default logger.
func DefaultLogger() *Logger {
	return defaultLogger
}
