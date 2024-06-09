// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package notifications

import (
	"context"
	"encoding/json"
	"time"
)

type Category uint8

const (
	Empty Category = iota
	Post
	Follow
	Like
	Comment
	Share
)

func (c Category) String() string {
	return [...]string{"", "Post", "Follow", "Like", "Comment", "Share"}[c]
}

func ToCategory(category string) Category {
	switch category {
	case "":
		return Empty
	case "Post":
		return Post
	case "Follow":
		return Follow
	case "Like":
		return Like
	case "Comment":
		return Comment
	case "Share":
		return Share
	default:
		return Empty
	}
}

func (c Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Category  Category  `json:"category"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Page struct {
	Total    uint64   `db:"total"              json:"total"`
	Offset   uint64   `db:"offset"             json:"offset"`
	Limit    uint64   `db:"limit"              json:"limit"`
	Category Category `db:"category,omitempty" json:"category,omitempty"`
	UserID   string   `db:"user_id,omitempty"  json:"user_id,omitempty"`
	IDs      []string `db:"ids,omitempty"      json:"ids,omitempty"`
	IsRead   *bool    `db:"is_read,omitempty"  json:"is_read,omitempty"`
}

type NotificationsPage struct {
	Page
	Notifications []Notification `json:"notifications"`
}

func (page NotificationsPage) MarshalJSON() ([]byte, error) {
	type Alias NotificationsPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Notifications == nil {
		a.Notifications = make([]Notification, 0)
	}

	return json.Marshal(a)
}

type Setting struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	IsEmailEnabled bool      `json:"is_email_enabled"`
	IsPushEnabled  bool      `json:"is_push_enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SettingsPage struct {
	Page
	Settings []Setting `json:"settings"`
}

func (page SettingsPage) MarshalJSON() ([]byte, error) {
	type Alias SettingsPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Settings == nil {
		a.Settings = make([]Setting, 0)
	}

	return json.Marshal(a)
}

//go:generate mockery --name Repository --output=./mocks --filename repository.go --quiet
type Repository interface {
	CreateNotification(ctx context.Context, notification Notification) (Notification, error)
	RetrieveNotification(ctx context.Context, id string) (Notification, error)
	RetrieveAllNotifications(ctx context.Context, page Page) (NotificationsPage, error)
	ReadNotification(ctx context.Context, userID, id string) error
	ReadAllNotifications(ctx context.Context, page Page) error
	DeleteNotification(ctx context.Context, id string) error
}

//go:generate mockery --name Service --output=./mocks --filename service.go --quiet
type Service interface {
	CreateNotification(ctx context.Context, notification Notification) (Notification, error)
	IdentifyUser(ctx context.Context, token string) (string, error)
	GetNewNotification(ctx context.Context, userID string) Notification
	RetrieveNotification(ctx context.Context, token string, id string) (Notification, error)
	RetrieveAllNotifications(ctx context.Context, token string, page Page) (NotificationsPage, error)
	ReadNotification(ctx context.Context, token string, id string) error
	ReadAllNotifications(ctx context.Context, token string, page Page) error
	DeleteNotification(ctx context.Context, token string, id string) error
}
