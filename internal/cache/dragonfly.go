// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Cacher = (*cache)(nil)

type cache struct {
	client   *redis.Client
	duration time.Duration
}

func NewCache(client *redis.Client, duration time.Duration) Cacher {
	return &cache{
		client:   client,
		duration: duration,
	}
}

func (c *cache) Add(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, string(data), c.duration).Err()
}

func (c *cache) Remove(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}

func (c *cache) Get(ctx context.Context, key string) interface{} {
	result, err := c.client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false
	}

	return result
}
