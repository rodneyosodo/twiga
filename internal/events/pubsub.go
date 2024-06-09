// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package events

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, msg map[string]interface{}) error
	Close() error
}

type EventHandler interface {
	Handle(ctx context.Context, msg map[string]interface{}) error
	Cancel() error
}

type SubscriberConfig struct {
	ID      string
	Topic   string
	Handler EventHandler
}

type Subscriber interface {
	Subscribe(ctx context.Context, cfg SubscriberConfig) error
	Unsubscribe(ctx context.Context, id, topic string) error
	Close() error
}

type PubSub interface {
	Publisher
	Subscriber
}
