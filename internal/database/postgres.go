package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DatabaseURL string

	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration

	ConnectTimeout time.Duration
	PingTimeout    time.Duration

	RetryAttempts int
}

func Connect(cfg Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL failed: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = cfg.MaxConns
	} else {
		poolConfig.MaxConns = 20
	}

	if cfg.MinConns > 0 {
		poolConfig.MinConns = cfg.MinConns
	}

	if cfg.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	} else {
		poolConfig.MaxConnLifetime = time.Hour
	}

	if cfg.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	} else {
		poolConfig.MaxConnIdleTime = 30 * time.Minute
	}

	connectTimeout := cfg.ConnectTimeout
	if connectTimeout == 0 {
		connectTimeout = 5 * time.Second
	}

	var pool *pgxpool.Pool

	for i := 0; i < max(1, cfg.RetryAttempts); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)

		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		cancel()

		if err == nil {
			break
		}

		log.Printf("DB connect attempt %d failed: %v", i+1, err)
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	pingTimeout := cfg.PingTimeout
	if pingTimeout == 0 {
		pingTimeout = 3 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("✅ PostgreSQL connected (pgx pool ready)")

	return pool, nil
}
