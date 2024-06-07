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
	General
	Follow
	Like
	Comment
)

func (c Category) String() string {
	return [...]string{"", "General", "Follow", "Like", "Comment"}[c]
}

func ToCategory(category string) Category {
	switch category {
	case "":
		return Empty
	case "General":
		return General
	case "Follow":
		return Follow
	case "Like":
		return Like
	case "Comment":
		return Comment
	default:
		return Empty
	}
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
	Total    uint64   `db:"total"    json:"total"`
	Offset   uint64   `db:"offset"   json:"offset"`
	Limit    uint64   `db:"limit"    json:"limit"`
	Category Category `db:"category" json:"category"`
	UserID   string   `db:"user_id"  json:"user_id"`
	IsRead   *bool    `db:"is_read"  json:"is_read"`
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
type Repository interface { //nolint:interfacebloat
	CreateNotification(ctx context.Context, notification Notification) (Notification, error)
	RetrieveNotification(ctx context.Context, id string) (Notification, error)
	RetrieveAllNotifications(ctx context.Context, page Page) (NotificationsPage, error)
	ReadNotification(ctx context.Context, userID, id string) error
	ReadAllNotifications(ctx context.Context, page Page) error
	DeleteNotification(ctx context.Context, id string) error

	CreateSetting(ctx context.Context, setting Setting) (Setting, error)
	RetrieveSetting(ctx context.Context, id string) (Setting, error)
	RetrieveAllSettings(ctx context.Context, page Page) (SettingsPage, error)
	UpdateSetting(ctx context.Context, setting Setting) error
	UpdateEmailSetting(ctx context.Context, id string, isEnabled bool) error
	UpdatePushSetting(ctx context.Context, id string, isEnabled bool) error
	DeleteSetting(ctx context.Context, id string) error
}

//go:generate mockery --name Service --output=./mocks --filename service.go --quiet
type Service interface { //nolint:interfacebloat
	CreateNotification(ctx context.Context, token string, notification Notification) (Notification, error)
	RetrieveNotification(ctx context.Context, token string, id string) (Notification, error)
	RetrieveAllNotifications(ctx context.Context, token string, page Page) (NotificationsPage, error)
	ReadNotification(ctx context.Context, token string, id string) error
	ReadAllNotifications(ctx context.Context, token string, page Page) error
	DeleteNotification(ctx context.Context, token string, id string) error

	CreateSetting(ctx context.Context, token string, setting Setting) (Setting, error)
	RetrieveSetting(ctx context.Context, token string, id string) (Setting, error)
	RetrieveAllSettings(ctx context.Context, token string, page Page) (SettingsPage, error)
	UpdateSetting(ctx context.Context, token string, setting Setting) error
	UpdateEmailSetting(ctx context.Context, token string, id string, isEnabled bool) error
	UpdatePushSetting(ctx context.Context, token string, id string, isEnabled bool) error
	DeleteSetting(ctx context.Context, token string, id string) error
}
