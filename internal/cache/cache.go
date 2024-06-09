// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package cache

import "context"

type Cacher interface {
	Add(ctx context.Context, key string, value interface{}) error
	Remove(ctx context.Context, key string) error
	Get(ctx context.Context, key string) interface{}
}
