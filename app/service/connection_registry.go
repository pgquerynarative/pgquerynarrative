package service

import (
	"strings"

	"github.com/pgquerynarrative/pgquerynarrative/app/catalog"
	"github.com/pgquerynarrative/pgquerynarrative/app/queryrunner"
)

type connectionResolver struct {
	defaultConnectionID string
	runners             map[string]*queryrunner.Runner
	loaders             map[string]*catalog.Loader
}

func newConnectionResolver(defaultID string, runners map[string]*queryrunner.Runner, loaders map[string]*catalog.Loader) connectionResolver {
	return connectionResolver{
		defaultConnectionID: defaultID,
		runners:             runners,
		loaders:             loaders,
	}
}

func (r connectionResolver) normalizedConnectionID(connectionID *string) string {
	if connectionID == nil || strings.TrimSpace(*connectionID) == "" {
		return r.defaultConnectionID
	}
	if _, ok := r.runners[strings.TrimSpace(*connectionID)]; ok {
		return strings.TrimSpace(*connectionID)
	}
	return r.defaultConnectionID
}

func (r connectionResolver) runnerFor(connectionID *string) *queryrunner.Runner {
	id := r.normalizedConnectionID(connectionID)
	return r.runners[id]
}

func (r connectionResolver) loaderFor(connectionID *string) *catalog.Loader {
	id := r.normalizedConnectionID(connectionID)
	return r.loaders[id]
}
