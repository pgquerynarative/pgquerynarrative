package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pgquerynarrative/pgquerynarrative/app/config"
	"github.com/pgquerynarrative/pgquerynarrative/app/errors"
)

type Pools struct {
	ReadOnly *pgxpool.Pool
	App      *pgxpool.Pool
}

func NewPools(ctx context.Context, cfg config.DatabaseConfig) (*Pools, error) {
	readOnlyURL := buildConnectionURL(
		cfg.ReadOnlyUser,
		cfg.ReadOnlyPassword,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	appURL := buildConnectionURL(
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	var readOnlyPool *pgxpool.Pool
	var err error
	maxRetries := 3
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		readOnlyPool, err = pgxpool.New(ctx, readOnlyURL)
		if err == nil {
			if pingErr := readOnlyPool.Ping(ctx); pingErr != nil {
				readOnlyPool.Close()
				err = fmt.Errorf("ping failed: %w", pingErr)
			} else {
				break
			}
		}
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2
		}
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrReadOnlyPoolFailed, err)
	}

	readOnlyPool.Config().MaxConns = int32(cfg.MaxConnections)
	readOnlyPool.Config().MaxConnLifetime = 30 * time.Minute
	readOnlyPool.Config().MinConns = 2

	var appPool *pgxpool.Pool
	retryDelay = 2 * time.Second
	for i := 0; i < maxRetries; i++ {
		appPool, err = pgxpool.New(ctx, appURL)
		if err == nil {
			if pingErr := appPool.Ping(ctx); pingErr != nil {
				appPool.Close()
				err = fmt.Errorf("ping failed: %w", pingErr)
			} else {
				break
			}
		}
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2
		}
	}
	if err != nil {
		readOnlyPool.Close()
		return nil, fmt.Errorf("%w: %v", errors.ErrAppPoolFailed, err)
	}

	appPool.Config().MaxConns = int32(cfg.MaxConnections)
	appPool.Config().MaxConnLifetime = 30 * time.Minute
	appPool.Config().MinConns = 2

	return &Pools{
		ReadOnly: readOnlyPool,
		App:      appPool,
	}, nil
}

func buildConnectionURL(user, password, host string, port int, database, sslMode string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		database,
		sslMode,
	)
}

func (p *Pools) Close() {
	if p == nil {
		return
	}
	if p.ReadOnly != nil {
		p.ReadOnly.Close()
	}
	if p.App != nil {
		p.App.Close()
	}
}

func (p *Pools) Health(ctx context.Context) error {
	if err := p.ReadOnly.Ping(ctx); err != nil {
		return fmt.Errorf("%w (read-only): %v", errors.ErrPoolHealthCheckFailed, err)
	}
	if err := p.App.Ping(ctx); err != nil {
		return fmt.Errorf("%w (app): %v", errors.ErrPoolHealthCheckFailed, err)
	}
	return nil
}
