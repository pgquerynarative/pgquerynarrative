package service

import (
	"context"
	"sort"

	"github.com/pgquerynarrative/pgquerynarrative/gen/connections"
)

// ConnectionsService exposes configured data connections (safe metadata only).
type ConnectionsService struct {
	items []*connections.ConnectionInfo
}

// NewConnectionsService creates a connections service from id/name map.
func NewConnectionsService(items []*connections.ConnectionInfo) *ConnectionsService {
	return &ConnectionsService{items: items}
}

// List returns all configured connection IDs and names.
func (s *ConnectionsService) List(context.Context) (*connections.ConnectionListResult, error) {
	out := make([]*connections.ConnectionInfo, len(s.items))
	copy(out, s.items)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return &connections.ConnectionListResult{Items: out}, nil
}
