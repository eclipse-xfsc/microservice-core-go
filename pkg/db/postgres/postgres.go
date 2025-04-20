package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sethvargo/go-retry"
)

// ConnectRetry tries to connect to the database given by Config. If the connection
// can't be established, it will retry periodically for maxDuration, using a fibonacci
// backoff starting with 500ms. It also uses a jitter of +/-50ms to avoid all instances
// trying to connect at the very same time. All errors not returned immediately
// (retryable errors) will be sent through errChan
func ConnectRetry(ctx context.Context, conf Config, maxDuration time.Duration, errChan chan<- error) (conn *pgxpool.Pool, err error) {
	if errChan != nil {
		defer close(errChan)
	}

	backoff := retry.NewFibonacci(500 * time.Millisecond)       // backoff 1s, 1s, 2s, 3s, 5s, ...
	backoff = retry.WithCappedDuration(time.Second*10, backoff) // backoff max 10s
	backoff = retry.WithJitter(50*time.Millisecond, backoff)    // 1s backoff, 50ms jitter -> random between 9950 and 10050ms
	backoff = retry.WithMaxDuration(maxDuration, backoff)       // stop retries after maxDuration

	err = retry.Do(ctx, backoff, func(ctx context.Context) error {
		conn, err = pgxpool.New(ctx, conf.DSN())
		if err != nil {
			return err
		}

		if err := conn.Ping(ctx); err != nil {
			if !isRetryableError(err) {
				return err
			}

			if errChan != nil {
				errChan <- err
			}

			return retry.RetryableError(err)
		}

		return nil
	})

	return
}

func MigrateUP(pool *pgxpool.Pool, migrations embed.FS, migrationsPath string) error {
	source, err := httpfs.New(http.FS(migrations), migrationsPath)
	if err != nil {
		return fmt.Errorf("could not initialize source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("httpfs", source, pool.Config().ConnString())
	if err != nil {
		return fmt.Errorf("failed to init migration: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func isRetryableError(err error) bool {
	if _, ok := err.(*pgconn.ConnectError); ok {
		return true
	}

	return false
}
