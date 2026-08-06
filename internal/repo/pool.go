package repo

import (
	"context"
	"log/slog"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns        = 20
	defaultMaxConnLifetime = 30 * time.Minute
	defaultMaxConnIdleTime = 5 * time.Minute
	defaultConnectTimeout  = 5 * time.Second
)

// Config holds the discrete connection parameters for the database, mirroring
// how each part (host, port, user, password, name, sslmode) is supplied
// independently by the deployment environment. The connection URL pgx needs
// is built from these, never accepted as a pre-built string, so a password
// containing special characters (@, :, /, %) is percent-encoded correctly
// instead of breaking a hand-spliced string.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string

	// Pool tuning — zero value means "use the package default" (see NewPool).
	MaxConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	ConnectTimeout  time.Duration
}

func (cfg Config) dsn() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   "/" + cfg.Name,
	}

	q := url.Values{}
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func NewPool(ctx context.Context, cfg Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.dsn())
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = orDefault(cfg.MaxConns, defaultMaxConns)
	poolCfg.MaxConnLifetime = orDefault(cfg.MaxConnLifetime, defaultMaxConnLifetime)
	poolCfg.MaxConnIdleTime = orDefault(cfg.MaxConnIdleTime, defaultMaxConnIdleTime)
	poolCfg.ConnConfig.ConnectTimeout = orDefault(cfg.ConnectTimeout, defaultConnectTimeout)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	logger.Info("connected to database", "host", cfg.Host, "name", cfg.Name)
	return pool, nil
}

// orDefault returns fallback when v is the zero value, otherwise v.
func orDefault[T comparable](v, fallback T) T {
	var zero T
	if v == zero {
		return fallback
	}
	return v
}
