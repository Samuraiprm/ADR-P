package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

type Option func(*config)

type config struct {
	addr     string
	password string
	db       int
}

func WithAddr(addr string) Option {
	return func(c *config) {
		if addr != "" {
			c.addr = addr
		}
	}
}

func WithPassword(password string) Option {
	return func(c *config) {
		c.password = password
	}
}

func WithDB(db int) Option {
	return func(c *config) {
		c.db = db
	}
}

func NewClient(opts ...Option) *Client {
	cfg := &config{
		addr:     "localhost:6379",
		password: "",
		db:       0,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.addr,
		Password:     cfg.password,
		DB:           cfg.db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	return &Client{rdb: rdb}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) IsConnected(ctx context.Context) bool {
	return c.rdb.Ping(ctx).Err() == nil
}

func (c *Client) XAdd(ctx context.Context, stream string, values map[string]interface{}) error {
	return c.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Err()
}

func (c *Client) Close() error {
	return c.rdb.Close()
}
