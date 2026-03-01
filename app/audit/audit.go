// Package audit writes security and API access events to app.audit_logs.
// Used by the HTTP server to record API_REQUEST, AUTH_FAILURE, AUTH_SUCCESS, RATE_LIMIT_EXCEEDED.
package audit

import (
	"context"
	"encoding/json"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Event types for audit_logs.event_type.
const (
	EventAPIRequest        = "API_REQUEST"
	EventAuthFailure       = "AUTH_FAILURE"
	EventAuthSuccess       = "AUTH_SUCCESS"
	EventRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	EventUnauthorized      = "UNAUTHORIZED_ACCESS"
)

// Entry represents a single audit log row.
type Entry struct {
	EventType  string
	EntityType string
	EntityID   *string
	Details    map[string]interface{}
	UserID     string
	IP         string
	UserAgent  string
}

// Store writes audit entries to the database.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore returns an audit store that writes to the given app pool.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// Record inserts one audit entry. It is best-effort; errors are not returned
// so that request handling is not failed by audit write failures.
func (s *Store) Record(ctx context.Context, e Entry) {
	if s == nil || s.pool == nil {
		return
	}
	detailsJSON, _ := json.Marshal(e.Details)
	var ip net.IP
	if e.IP != "" {
		ip = net.ParseIP(e.IP)
	}
	_, _ = s.pool.Exec(ctx,
		`INSERT INTO app.audit_logs (event_type, entity_type, entity_id, details, user_id, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, NULLIF($5,''), $6, NULLIF($7,''))`,
		e.EventType, e.EntityType, e.EntityID, detailsJSON, e.UserID, ip, e.UserAgent,
	)
}
