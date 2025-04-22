package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"plata/internal/common/log"
	"plata/internal/config"
	"time"
)

type DBGetter interface {
	Primary() *sqlx.DB
	Replica() *sqlx.DB
}

type PgDB struct {
	prime   *sqlx.DB
	replica *sqlx.DB
	log     log.Logger
}

func (pg *PgDB) Primary() *sqlx.DB {
	return pg.prime
}

func (pg *PgDB) Replica() *sqlx.DB {
	if pg.replica != nil {
		return pg.replica
	}
	return pg.prime
}

func NewPostgresDB(cfg config.PostgresConfig, logger log.Logger) (*PgDB, error) {
	dsnPrimary := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.HostPrimary,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	prime, err := sqlx.Connect("postgres", dsnPrimary)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to primary: %w", err)
	}
	if err = prime.Ping(); err != nil {
		return nil, fmt.Errorf("ping to primary failed: %w", err)
	}

	idleTime, err := time.ParseDuration(cfg.ConnMaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("invalid ConnMaxIdleTime format: %w", err)
	}
	prime.SetMaxOpenConns(cfg.MaxOpenConns)
	prime.SetMaxIdleConns(cfg.MaxIdleConns)
	prime.SetConnMaxIdleTime(idleTime)

	var replica *sqlx.DB
	if cfg.HostReplica != "" {
		dsnReplica := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.HostReplica,
			cfg.Port,
			cfg.Username,
			cfg.Password,
			cfg.Database,
			cfg.SSLMode,
		)

		replica, err = sqlx.Connect("postgres", dsnReplica)
		if err != nil {
			logger.Warnf("Failed to connect to replica: %v", err)
			replica = nil
		} else {
			logger.Infof("Connected to replica")
		}
	}

	return &PgDB{
		prime:   prime,
		replica: replica,
		log:     logger,
	}, nil
}

func (pg *PgDB) Stop() {
	if err := pg.prime.Close(); err != nil {
		pg.log.Errorf("Failed to close primary DB: %v", err)
	}
	if pg.replica != nil {
		if err := pg.replica.Close(); err != nil {
			pg.log.Errorf("Failed to close replica DB: %v", err)
		}
	}
}
