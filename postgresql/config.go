package postgresql

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

// Config struct for postgresql config
type Config struct {
	Port              int           `mapstructure:"port"`
	Host              string        `mapstructure:"host"`
	User              string        `mapstructure:"user"`
	Database          string        `mapstructure:"database"`
	Password          string        `mapstructure:"password"`
	Secured           bool          `mapstructure:"secured"`
	MaxConnections    int32         `mapstructure:"max_connections"`
	MinConnections    int32         `mapstructure:"min_connections"`
	MaxConnectionAge  time.Duration `mapstructure:"max_connection_age"`
	AcquireTimeout    time.Duration `mapstructure:"acquire_timeout"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period"`
	LoggerEnabled     bool          `mapstructure:"logger_enabled"`
	LogLevel          pgx.LogLevel  `mapstructure:"log_level"`
	KeepAlive         time.Duration `mapstructure:"keep_alive"`
	Schema            string        `mapstructure:"schema"`
}

// ConnString return connection string
func (c Config) ConnString() string {
	return c.ConnStringFor(c.Host)
}

func (c Config) ConnStringFor(host string) string {
	sslmode := "disable"
	if c.Secured {
		sslmode = "require"
	}
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s sslmode=%s",
		host, c.Port, c.User, sslmode,
	)
	if c.Database != "" {
		connStr = connStr + " dbname=" + c.Database
	}
	if c.Password != "" {
		connStr = connStr + " password=" + c.Password
	}
	if c.Schema != "" {
		connStr = connStr + " search_path=" + c.Schema
	}
	return connStr
}

type wrappedDialer struct {
	*net.Dialer
}

func (d *wrappedDialer) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	start := time.Now()
	defer func() {
		l := log.Debug().
			Str("component", "pgx").
			Dur("duration", time.Since(start)).
			Str("network", network).
			Str("address", address)
		if err != nil {
			l = l.Err(err)
		}
		l.Msg("dial postgres")
	}()
	conn, err = d.Dialer.DialContext(ctx, network, address)
	return
}

type wrappedResolver struct {
	*net.Resolver
}

func (r *wrappedResolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
	start := time.Now()
	defer func() {
		l := log.Debug().
			Str("component", "pgx").
			Dur("duration", time.Since(start)).
			Str("host", host).
			Strs("resolved_addrs", addrs)
		if err != nil {
			l = l.Err(err)
		}
		l.Msg("resolve postgres host")
	}()
	addrs, err = r.Resolver.LookupHost(ctx, host)
	return
}

type ConnectionPoolOption func(*ConnectionPool) error

// NewConnectionPool return new Connection Pool
func NewConnectionPool(conf Config, opts ...ConnectionPoolOption) (DB, error) {
	poolConfig, err := pgxpool.ParseConfig(conf.ConnString())
	if err != nil {
		return nil, err
	}

	conf.ApplyPoolConfig(poolConfig)

	p, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	connectionPool := &ConnectionPool{Pool: p}

	for _, opt := range opts {
		err := opt(connectionPool)
		if err != nil {
			return nil, err
		}
	}

	return connectionPool, nil

}

func (conf Config) ApplyPoolConfig(poolConfig *pgxpool.Config) {
	poolConfig.MaxConns = conf.MaxConnections
	poolConfig.MinConns = conf.MinConnections
	poolConfig.MaxConnLifetime = conf.MaxConnectionAge
	poolConfig.HealthCheckPeriod = conf.HealthCheckPeriod
	if conf.LoggerEnabled {
		dialer := &wrappedDialer{
			&net.Dialer{
				KeepAlive: conf.KeepAlive,
			},
		}
		resolver := &wrappedResolver{
			net.DefaultResolver,
		}
		poolConfig.ConnConfig.Logger = zerologadapter.NewLogger(log.Logger)
		poolConfig.ConnConfig.LogLevel = conf.LogLevel
		poolConfig.ConnConfig.DialFunc = dialer.DialContext
		poolConfig.ConnConfig.LookupFunc = resolver.LookupHost
	}
}

func ReplicaPoolsFromConfig(replicaConfigs ...Config) ([]*pgxpool.Pool, error) {
	var replicas []*pgxpool.Pool
	for _, conf := range replicaConfigs {
		replicaPoolConfig, err := pgxpool.ParseConfig(conf.ConnString())
		if err != nil {
			return nil, err
		}

		conf.ApplyPoolConfig(replicaPoolConfig)

		replica, err := pgxpool.ConnectConfig(context.Background(), replicaPoolConfig)
		if err != nil {
			return nil, err
		}

		replicas = append(replicas, replica)
	}

	return replicas, nil
}
