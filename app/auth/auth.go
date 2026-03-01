// Package auth provides API key (Bearer token) validation for the HTTP server.
// When SECURITY_AUTH_ENABLED is true, API and web export routes require
// Authorization: Bearer <SECURITY_API_KEY>.
package auth

import (
	"net/http"
	"strings"
)

// ContextKey is the type for auth identity in request context.
type ContextKey string

const (
	// IdentityContextKey is the context key for the authenticated identity (e.g. "api-key" or "").
	IdentityContextKey ContextKey = "auth_identity"
)

// ValidateRequest checks the request for a valid Bearer token when expectedAPIKey is non-empty.
// It returns the identity to use for audit (e.g. "api-key") and true if the request is allowed.
// When expectedAPIKey is empty, all requests are allowed and identity is "".
func ValidateRequest(r *http.Request, expectedAPIKey string) (identity string, ok bool) {
	if expectedAPIKey == "" {
		return "", true
	}
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	if token == "" {
		return "", false
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(token, prefix) {
		return "", false
	}
	token = strings.TrimSpace(token[len(prefix):])
	if token == "" || token != expectedAPIKey {
		return "", false
	}
	return "api-key", true
}
