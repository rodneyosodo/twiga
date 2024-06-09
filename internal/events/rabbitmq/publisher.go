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
	"fmt"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rodneyosodo/twiga/internal/events"
)

var _ events.Publisher = (*publisher)(nil)

type publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	prefix   string
	exchange string
}

func NewPublisher(url string) (events.Publisher, error) {
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

	ret := &publisher{
		conn:     conn,
		channel:  ch,
		prefix:   chansPrefix,
		exchange: exchangeName,
	}

	return ret, nil
}

func (pub *publisher) Publish(ctx context.Context, topic string, msg map[string]interface{}) error {
	if topic == "" {
		return ErrEmptyTopic
	}
	msg["timestamp"] = time.Now().UnixNano()
	msg["topic"] = topic

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s.%s", pub.prefix, topic)

	subject = formatTopic(subject)

	err = pub.channel.PublishWithContext(
		ctx,
		pub.exchange,
		subject,
		false,
		false,
		amqp.Publishing{
			Headers:     amqp.Table{},
			ContentType: "application/octet-stream",
			AppId:       "twiga-publisher",
			Body:        data,
		})
	if err != nil {
		return err
	}

	return nil
}

func (pub *publisher) Close() error {
	return pub.conn.Close()
}

func formatTopic(topic string) string {
	return strings.ReplaceAll(topic, ">", "#")
}
