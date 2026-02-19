// Package catalog provides read-only database schema discovery from
// PostgreSQL information_schema for allowed schemas (e.g. demo).
package catalog

import (
	"context"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"

	schema "github.com/pgquerynarrative/pgquerynarrative/gen/schema"
)

// Loader loads schema metadata from the database using the read-only pool.
// Only tables in allowed schemas are returned, so the result matches
// what the query validator permits.
type Loader struct {
	pool           *pgxpool.Pool
	allowedSchemas []string
}

// NewLoader creates a catalog loader that queries information_schema
// and returns only the given allowed schema names (e.g. []string{"demo"}).
func NewLoader(pool *pgxpool.Pool, allowedSchemas []string) *Loader {
	return &Loader{pool: pool, allowedSchemas: allowedSchemas}
}

// Load returns the list of allowed schemas with their tables and columns.
// It uses the read-only pool so only objects visible to that user are included.
func (l *Loader) Load(ctx context.Context) (*schema.SchemaResult, error) {
	if len(l.allowedSchemas) == 0 {
		return &schema.SchemaResult{Schemas: []*schema.SchemaInfo{}}, nil
	}

	// Build schema name list for SQL IN clause. We use a safe allowlist.
	// Columns: table_schema, table_name, column_name, data_type, ordinal_position
	const q = `
		SELECT table_schema, table_name, column_name, data_type, ordinal_position
		FROM information_schema.columns
		WHERE table_schema = ANY($1)
		  AND table_catalog = current_database()
		ORDER BY table_schema, table_name, ordinal_position
	`
	rows, err := l.pool.Query(ctx, q, l.allowedSchemas)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type row struct {
		tableSchema string
		tableName   string
		columnName  string
		dataType    string
		ordPosition int32
	}
	var raw []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.tableSchema, &r.tableName, &r.columnName, &r.dataType, &r.ordPosition); err != nil {
			return nil, err
		}
		raw = append(raw, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Group by schema -> table -> columns
	schemaMap := make(map[string]map[string][]*schema.ColumnInfo)
	for _, r := range raw {
		if schemaMap[r.tableSchema] == nil {
			schemaMap[r.tableSchema] = make(map[string][]*schema.ColumnInfo)
		}
		tables := schemaMap[r.tableSchema]
		tables[r.tableName] = append(tables[r.tableName], &schema.ColumnInfo{Name: r.columnName, Type: r.dataType})
	}

	var schemas []*schema.SchemaInfo
	for _, name := range l.allowedSchemas {
		tablesMap, ok := schemaMap[name]
		if !ok {
			schemas = append(schemas, &schema.SchemaInfo{Name: name, Tables: []*schema.TableInfo{}})
			continue
		}
		var tables []*schema.TableInfo
		for tableName, cols := range tablesMap {
			tables = append(tables, &schema.TableInfo{Name: tableName, Columns: cols})
		}
		sort.Slice(tables, func(i, j int) bool { return tables[i].Name < tables[j].Name })
		schemas = append(schemas, &schema.SchemaInfo{Name: name, Tables: tables})
	}
	return &schema.SchemaResult{Schemas: schemas}, nil
}
