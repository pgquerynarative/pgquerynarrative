// Package queryrunner provides SQL query validation to ensure queries are safe
// and read-only. It prevents SQL injection and enforces security policies.
package queryrunner

import (
	"regexp"
	"strings"

	"github.com/pgquerynarrative/pgquerynarrative/app/errors"
)

// Validator validates SQL queries to ensure they are safe to execute.
// It enforces:
//   - Read-only queries (SELECT/WITH only)
//   - Single statement execution
//   - Allowed schemas only
//   - No dangerous keywords
type Validator struct {
	allowedSchemas   map[string]bool // Set of allowed schema names (lowercase)
	maxQueryLength   int             // Maximum query length in bytes
	disallowedTokens []string        // Keywords that are not allowed
}

// NewValidator creates a new query validator with the specified configuration.
//
// Parameters:
//   - allowedSchemas: List of schema names that queries are allowed to access
//   - maxQueryLength: Maximum query length in bytes (prevents DoS)
//
// Returns a configured Validator instance.
func NewValidator(allowedSchemas []string, maxQueryLength int) *Validator {
	// Build schema map for O(1) lookup
	schemaMap := make(map[string]bool, len(allowedSchemas))
	for _, schema := range allowedSchemas {
		schemaMap[strings.ToLower(schema)] = true
	}

	return &Validator{
		allowedSchemas: schemaMap,
		maxQueryLength: maxQueryLength,
		disallowedTokens: []string{
			// Data modification
			"insert", "update", "delete", "truncate",
			// Schema modification
			"drop", "create", "alter",
			// Privilege modification
			"grant", "revoke",
			// Execution
			"execute",
			// Data export
			"copy",
			// Cursors and transactions
			"declare", "cursor",
			// Performance analysis (can be slow)
			"explain analyze",
			// Maintenance (can be slow or lock tables)
			"vacuum", "analyze", "reindex", "cluster",
		},
	}
}

// Validate checks if a SQL query is safe to execute.
//
// Validation rules:
//   - Query length must not exceed maxQueryLength
//   - Must be a single statement (no semicolons except trailing)
//   - Must start with SELECT or WITH (read-only)
//   - Must not contain disallowed keywords
//   - Must only reference allowed schemas
//
// Parameters:
//   - sql: SQL query string to validate
//
// Returns:
//   - nil if query is valid
//   - Error describing why the query is invalid
func (v *Validator) Validate(sql string) error {
	// Check query length
	if len(sql) > v.maxQueryLength {
		return errors.ErrQueryTooLong
	}

	// Normalize query: trim whitespace and remove trailing semicolon
	trimmed := strings.TrimSpace(sql)
	trimmed = strings.TrimSuffix(trimmed, ";")
	trimmed = strings.TrimSpace(trimmed)

	// Check for multiple statements (semicolons in the middle)
	if strings.Contains(trimmed, ";") {
		return errors.ErrMultipleStatements
	}

	// Convert to lowercase for keyword checking
	lower := strings.ToLower(trimmed)

	// Must start with SELECT or WITH (read-only queries only)
	if !strings.HasPrefix(lower, "select") && !strings.HasPrefix(lower, "with") {
		return errors.ErrOnlySelectAllowed
	}

	// Check for disallowed keywords
	for _, token := range v.disallowedTokens {
		if strings.Contains(lower, token) {
			return errors.ErrDisallowedKeyword
		}
	}

	// Check schema access (if schema restrictions are configured)
	if len(v.allowedSchemas) > 0 {
		// Match schema.table only in FROM/JOIN (avoids false positives from alias.column)
		schemaPattern := regexp.MustCompile(`(?:from|join)\s+([a-z_][a-z0-9_]*)\.`)
		matches := schemaPattern.FindAllStringSubmatch(lower, -1)

		for _, match := range matches {
			if len(match) > 1 {
				schemaName := match[1]
				if !v.allowedSchemas[schemaName] {
					return errors.ErrSchemaNotAllowed
				}
			}
		}
	}

	return nil
}
