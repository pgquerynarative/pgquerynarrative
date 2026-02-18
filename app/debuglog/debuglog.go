// Package debuglog provides optional verbose logging when LOG_DEBUG is enabled.
package debuglog

import (
	"log"
	"os"
	"strings"
)

func enabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_DEBUG")))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

// Log writes a log message when LOG_DEBUG is set (1, true, yes, on).
func Log(format string, args ...interface{}) {
	if enabled() {
		log.Printf("[pgquerynarrative] "+format, args...)
	}
}
