// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package consumer

import (
	"context"

	"github.com/rodneyosodo/twiga/internal/events"
	"github.com/rodneyosodo/twiga/notifications"
)

const (
	postCreated = "posts.created"
	postUpdated = "posts.updated"
	postDeleted = "posts.deleted"

	commentCreated = "comments.created"
	commentUpdated = "comments.updated"
	commentDeleted = "comments.deleted"

	likeCreated = "likes.created"
	likeDeleted = "likes.deleted"

	shareCreated = "shares.created"
	shareDeleted = "shares.deleted"

	followCreated = "followers.created"
)

type eventHandler struct {
	notifications.Service
}

func NewEventHandler(svc notifications.Service) events.EventHandler {
	return &eventHandler{Service: svc}
}

func (eh *eventHandler) Handle(ctx context.Context, event map[string]interface{}) error {
	var notification notifications.Notification

	switch event["topic"] {
	case postCreated, postUpdated:
		notification = decodeEntityCreatedEvent(event, notifications.Post)
	case postDeleted:
		notification = decodeEntityDeletedEvent(event, notifications.Post)
	case commentCreated, commentUpdated:
		notification = decodeEntityCreatedEvent(event, notifications.Comment)
	case commentDeleted:
		notification = decodeEntityDeletedEvent(event, notifications.Comment)
	case likeCreated:
		notification = decodeEntityCreatedEvent(event, notifications.Like)
	case likeDeleted:
		notification = decodeEntityDeletedEvent(event, notifications.Like)
	case shareCreated:
		notification = decodeEntityCreatedEvent(event, notifications.Share)
	case shareDeleted:
		notification = decodeEntityDeletedEvent(event, notifications.Share)
	case followCreated:
		notification = decodeFollowCreatedEvent(event)
	}

	if _, err := eh.CreateNotification(ctx, notification); err != nil {
		return err
	}

	return nil
}

func (eh *eventHandler) Cancel() error {
	return nil
}

func decodeEntityCreatedEvent(event map[string]interface{}, category notifications.Category) notifications.Notification {
	return notifications.Notification{
		UserID:   read(event, "user_id", ""),
		Category: category,
		Content:  read(event, "content", ""),
	}
}

func decodeEntityDeletedEvent(event map[string]interface{}, category notifications.Category) notifications.Notification {
	return notifications.Notification{
		UserID:   read(event, "user_id", ""),
		Category: category,
		Content:  read(event, "id", ""),
	}
}

func decodeFollowCreatedEvent(event map[string]interface{}) notifications.Notification {
	return notifications.Notification{
		UserID:   read(event, "follower_id", ""),
		Category: notifications.Follow,
		Content:  read(event, "followee_id", ""),
	}
}

func read(event map[string]interface{}, key, def string) string {
	val, ok := event[key].(string)
	if !ok {
		return def
	}

	return val
}
