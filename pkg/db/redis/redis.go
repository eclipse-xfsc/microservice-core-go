package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-retry"
)

type Client struct {
	Rdb        redis.Cmdable
	DefaultTTL time.Duration
}

func ConnectRetry(ctx context.Context, config Config, defaultTTL time.Duration, errChan chan<- error) (*Client, error) {
	var rdb redis.Cmdable
	if config.IsCluster {
		nodes := strings.Split(config.Hosts, ";")
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          nodes,
			Username:       config.Username,
			Password:       config.Password,
			RouteByLatency: true,
			DialTimeout:    10 * time.Second,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxRedirects:   10,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:         config.Hosts,
			Username:     config.Username,
			Password:     config.Password,
			DB:           config.Database,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		})
	}

	backoff := retry.NewFibonacci(time.Millisecond * 500)
	backoff = retry.WithCappedDuration(time.Second*30, backoff)
	backoff = retry.WithJitter(time.Millisecond*50, backoff)
	backoff = retry.WithMaxDuration(time.Minute, backoff)

	if err := retry.Do(ctx, backoff, func(ctx context.Context) error {
		if _, err := rdb.Ping(ctx).Result(); err != nil {
			err = fmt.Errorf("connection check failed: error: %w", err)
			errChan <- err

			return retry.RetryableError(err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Client{
		Rdb:        rdb,
		DefaultTTL: defaultTTL,
	}, nil

}
