// Package queryrunner provides SQL query validation to ensure queries are safe
// and read-only. It prevents SQL injection and enforces security policies.
package queryrunner

import (
	"strings"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/pgquerynarrative/pgquerynarrative/app/errors"
)

// Validator validates SQL queries to ensure they are safe to execute.
// It enforces:
//   - Read-only queries (SELECT/WITH only)
//   - Single statement execution
//   - Allowed schemas only
//   - No dangerous keywords
type Validator struct {
	allowedSchemas map[string]bool // Set of allowed schema names (lowercase)
	maxQueryLength int             // Maximum query length in bytes
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
	}
}

// Validate checks if a SQL query is safe to execute.
//
// Validation rules:
//   - Query length must not exceed maxQueryLength
//   - Must be a single statement
//   - Must be a top-level SELECT/WITH SELECT
//   - Must reject write/unsafe statements in nested CTE/subquery contexts
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

	// Normalize query: trim whitespace
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return errors.ErrOnlySelectAllowed
	}

	stmts, err := parser.Parse(trimmed)
	if err != nil {
		return errors.ErrOnlySelectAllowed
	}
	if len(stmts) != 1 {
		return errors.ErrMultipleStatements
	}
	if stmts[0].AST == nil {
		return errors.ErrOnlySelectAllowed
	}

	if err := v.validateReadOnlyStatement(stmts[0].AST, true); err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateReadOnlyStatement(stmt tree.Statement, topLevel bool) error {
	switch s := stmt.(type) {
	case *tree.Select:
		return v.validateSelect(s, topLevel)
	case *tree.ParenSelect:
		if s.Select == nil {
			return errors.ErrOnlySelectAllowed
		}
		return v.validateSelect(s.Select, topLevel)
	default:
		if topLevel {
			return errors.ErrOnlySelectAllowed
		}
		return errors.ErrDisallowedKeyword
	}
}

func (v *Validator) validateSelect(stmt *tree.Select, topLevel bool) error {
	if stmt == nil {
		if topLevel {
			return errors.ErrOnlySelectAllowed
		}
		return errors.ErrDisallowedKeyword
	}

	if len(stmt.Locking) > 0 {
		return errors.ErrDisallowedKeyword
	}

	if stmt.With != nil {
		for _, cte := range stmt.With.CTEList {
			if cte == nil || cte.Stmt == nil {
				return errors.ErrDisallowedKeyword
			}
			if err := v.validateReadOnlyStatement(cte.Stmt, false); err != nil {
				return err
			}
		}
	}

	return v.validateSelectStatement(stmt.Select, topLevel)
}

func (v *Validator) validateSelectStatement(stmt tree.SelectStatement, topLevel bool) error {
	switch s := stmt.(type) {
	case *tree.SelectClause:
		for _, tableExpr := range s.From.Tables {
			if err := v.validateTableExpr(tableExpr); err != nil {
				return err
			}
		}
		return nil
	case *tree.UnionClause:
		if s.Left != nil {
			if err := v.validateSelect(s.Left, false); err != nil {
				return err
			}
		}
		if s.Right != nil {
			if err := v.validateSelect(s.Right, false); err != nil {
				return err
			}
		}
		return nil
	case *tree.ParenSelect:
		if s.Select == nil {
			return errors.ErrDisallowedKeyword
		}
		return v.validateSelect(s.Select, false)
	case *tree.ValuesClause:
		// Keep the previous policy: only SELECT/WITH-style statements.
		if topLevel {
			return errors.ErrOnlySelectAllowed
		}
		return nil
	default:
		if topLevel {
			return errors.ErrOnlySelectAllowed
		}
		return errors.ErrDisallowedKeyword
	}
}

func (v *Validator) validateTableExpr(expr tree.TableExpr) error {
	switch t := expr.(type) {
	case *tree.AliasedTableExpr:
		return v.validateTableExpr(t.Expr)
	case *tree.ParenTableExpr:
		return v.validateTableExpr(t.Expr)
	case *tree.JoinTableExpr:
		if err := v.validateTableExpr(t.Left); err != nil {
			return err
		}
		return v.validateTableExpr(t.Right)
	case *tree.Subquery:
		if t.Select == nil {
			return errors.ErrDisallowedKeyword
		}
		return v.validateSelectStatement(t.Select, false)
	case *tree.StatementSource:
		if t.Statement == nil {
			return errors.ErrDisallowedKeyword
		}
		return v.validateReadOnlyStatement(t.Statement, false)
	case *tree.TableName:
		return v.validateSchemaName(t.Schema(), t.ExplicitSchema)
	default:
		return nil
	}
}

func (v *Validator) validateSchemaName(schemaName string, isExplicit bool) error {
	if !isExplicit || len(v.allowedSchemas) == 0 {
		return nil
	}
	schema := strings.ToLower(strings.TrimSpace(schemaName))
	if schema == "" {
		return nil
	}
	if !v.allowedSchemas[schema] {
		return errors.ErrSchemaNotAllowed
	}
	return nil
}
