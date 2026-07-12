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

func (c *Client) XReadGroup(ctx context.Context, group, consumer string, streams map[string]string, count int64, block time.Duration) ([]redis.XStream, error) {
	streamKeys := make([]string, 0, len(streams)*2)
	streamIDs := make([]string, 0, len(streams))
	for k, v := range streams {
		streamKeys = append(streamKeys, k)
		streamIDs = append(streamIDs, v)
	}
	streamArgs := append(streamKeys, streamIDs...)

	return c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  streamArgs,
		Count:    count,
		Block:    block,
	}).Result()
}

func (c *Client) XAck(ctx context.Context, stream, group string, ids ...string) error {
	return c.rdb.XAck(ctx, stream, group, ids...).Err()
}

func (c *Client) XGroupCreate(ctx context.Context, stream, group, startID string) error {
	return c.rdb.XGroupCreateMkStream(ctx, stream, group, startID).Err()
}

func (c *Client) Close() error {
	return c.rdb.Close()
}
