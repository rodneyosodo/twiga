// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rodneyosodo/twiga/internal/events"
)

const (
	SubjectAllEvents = "events.#"

	exchangeName = "events"
	chansPrefix  = "events"
)

var (
	ErrNotSubscribed = errors.New("not subscribed")
	ErrEmptyTopic    = errors.New("empty topic")
	ErrEmptyID       = errors.New("empty ID")
)

var _ events.PubSub = (*pubsub)(nil)

type subscription struct {
	cancel func() error
}

type pubsub struct {
	publisher
	logger        *slog.Logger
	subscriptions map[string]map[string]subscription
	mu            sync.Mutex
}

func NewPubSub(url string, logger *slog.Logger) (events.PubSub, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err := ch.ExchangeDeclare(exchangeName, amqp.ExchangeTopic, true, false, false, false, nil); err != nil {
		return nil, err
	}

	ret := &pubsub{
		publisher: publisher{
			conn:     conn,
			channel:  ch,
			exchange: exchangeName,
			prefix:   chansPrefix,
		},
		logger:        logger,
		subscriptions: make(map[string]map[string]subscription),
	}

	return ret, nil
}

func (ps *pubsub) Subscribe(ctx context.Context, cfg events.SubscriberConfig) error {
	if cfg.ID == "" {
		return ErrEmptyID
	}
	if cfg.Topic == "" {
		return ErrEmptyTopic
	}
	ps.mu.Lock()

	cfg.Topic = formatTopic(cfg.Topic)
	s, ok := ps.subscriptions[cfg.Topic]
	if ok {
		if _, ok := s[cfg.ID]; ok {
			ps.mu.Unlock()
			if err := ps.Unsubscribe(ctx, cfg.ID, cfg.Topic); err != nil {
				return err
			}

			ps.mu.Lock()
			s = ps.subscriptions[cfg.Topic]
		}
	}
	defer ps.mu.Unlock()
	if s == nil {
		s = make(map[string]subscription)
		ps.subscriptions[cfg.Topic] = s
	}

	clientID := fmt.Sprintf("%s-%s", cfg.Topic, cfg.ID)

	queue, err := ps.channel.QueueDeclare(clientID, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := ps.channel.QueueBind(queue.Name, cfg.Topic, ps.exchange, false, nil); err != nil {
		return err
	}

	msgs, err := ps.channel.Consume(queue.Name, clientID, true, false, false, false, nil)
	if err != nil {
		return err
	}
	go ps.handle(ctx, msgs, cfg.Handler)
	s[cfg.ID] = subscription{
		cancel: func() error {
			if err := ps.channel.Cancel(clientID, false); err != nil {
				return err
			}

			return cfg.Handler.Cancel()
		},
	}

	return nil
}

func (ps *pubsub) Unsubscribe(ctx context.Context, id, topic string) error {
	if id == "" {
		return ErrEmptyID
	}
	if topic == "" {
		return ErrEmptyTopic
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()

	topic = formatTopic(topic)
	s, ok := ps.subscriptions[topic]
	if !ok {
		return ErrNotSubscribed
	}

	current, ok := s[id]
	if !ok {
		return ErrNotSubscribed
	}
	if current.cancel != nil {
		if err := current.cancel(); err != nil {
			return err
		}
	}
	if err := ps.channel.QueueUnbind(topic, topic, exchangeName, nil); err != nil {
		return err
	}

	delete(s, id)
	if len(s) == 0 {
		delete(ps.subscriptions, topic)
	}

	return nil
}

func (ps *pubsub) handle(ctx context.Context, deliveries <-chan amqp.Delivery, h events.EventHandler) {
	for d := range deliveries {
		var msg map[string]interface{}
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			ps.logger.Warn(fmt.Sprintf("Failed to unmarshal received event: %s", err))

			return
		}
		if err := h.Handle(ctx, msg); err != nil {
			ps.logger.Warn(fmt.Sprintf("Failed to handle twiga event: %s", err))

			return
		}
	}
}
