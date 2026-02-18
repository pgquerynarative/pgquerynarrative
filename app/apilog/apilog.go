// Package apilog provides console logging for API errors and important events.
// All messages are written to stdout with a consistent prefix for visibility
// when running the server (e.g. validation failures, LLM errors, request errors).
package apilog

import (
	"log"
	"os"
)

var defaultLogger = log.New(os.Stdout, "[pgquerynarrative] ", log.LstdFlags)

// ValidationError logs a validation error (query or report) for the given endpoint.
func ValidationError(endpoint, name, message string) {
	defaultLogger.Printf("validation_error %s name=%s message=%s", endpoint, name, message)
}

// LLMError logs an LLM/narrative generation error.
func LLMError(message string) {
	defaultLogger.Printf("llm_error message=%s", message)
}

// APIError logs a generic API error response (e.g. 4xx/5xx).
func APIError(method, path string, status int) {
	defaultLogger.Printf("api_error %s %s status=%d", method, path, status)
}

// Request logs a successful API request when LOG_DEBUG is set (handled via debuglog for consistency).
// Use for high-value events like "report generated" if desired; for now errors are the focus.
func Request(endpoint, detail string) {
	// Same prefix so it appears in console with other logs
	defaultLogger.Printf("api %s %s", endpoint, detail)
}
