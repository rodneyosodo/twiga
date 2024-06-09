// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package producer

import (
	"context"
	"encoding/json"

	"github.com/rodneyosodo/twiga/internal/events"
	"github.com/rodneyosodo/twiga/users"
)

var _ users.Service = (*eventStore)(nil)

type eventStore struct {
	publisher events.Publisher
	svc       users.Service
}

func NewEventStore(pub events.Publisher, svc users.Service) users.Service {
	return &eventStore{
		publisher: pub,
		svc:       svc,
	}
}

func (e *eventStore) IssueToken(ctx context.Context, user users.User) (string, error) {
	return e.svc.IssueToken(ctx, user)
}

func (e *eventStore) RefreshToken(ctx context.Context, token string) (string, error) {
	return e.svc.RefreshToken(ctx, token)
}

func (e *eventStore) IdentifyUser(ctx context.Context, token string) (string, error) {
	return e.svc.IdentifyUser(ctx, token)
}

func (e *eventStore) CreateUser(ctx context.Context, user users.User) (users.User, error) {
	return e.svc.CreateUser(ctx, user)
}

func (e *eventStore) GetUserByID(ctx context.Context, token string, id string) (users.User, error) {
	return e.svc.GetUserByID(ctx, token, id)
}

func (e *eventStore) GetUsers(ctx context.Context, token string, page users.Page) (users.UsersPage, error) {
	return e.svc.GetUsers(ctx, token, page)
}

func (e *eventStore) UpdateUser(ctx context.Context, token string, user users.User) (users.User, error) {
	return e.svc.UpdateUser(ctx, token, user)
}

func (e *eventStore) UpdateUserUsername(ctx context.Context, token string, user users.User) (users.User, error) {
	return e.svc.UpdateUserUsername(ctx, token, user)
}

func (e *eventStore) UpdateUserPassword(ctx context.Context, token string, oldPassword, currentPassowrd string) error {
	return e.svc.UpdateUserPassword(ctx, token, oldPassword, currentPassowrd)
}

func (e *eventStore) UpdateUserEmail(ctx context.Context, token string, user users.User) (users.User, error) {
	return e.svc.UpdateUserEmail(ctx, token, user)
}

func (e *eventStore) UpdateUserBio(ctx context.Context, token string, user users.User) (users.User, error) {
	return e.svc.UpdateUserBio(ctx, token, user)
}

func (e *eventStore) UpdateUserPictureURL(ctx context.Context, token string, user users.User) (users.User, error) {
	return e.svc.UpdateUserPictureURL(ctx, token, user)
}

func (e *eventStore) UpdateUserPreferences(ctx context.Context, token string, user users.User) (users.User, error) {
	return e.svc.UpdateUserPreferences(ctx, token, user)
}

func (e *eventStore) DeleteUser(ctx context.Context, token string, id string) error {
	return e.svc.DeleteUser(ctx, token, id)
}

func (e *eventStore) CreatePreferences(ctx context.Context, token string, preference users.Preference) (users.Preference, error) {
	return e.svc.CreatePreferences(ctx, token, preference)
}

func (e *eventStore) GetPreferencesByUserID(ctx context.Context, token, id string) (users.Preference, error) {
	return e.svc.GetPreferencesByUserID(ctx, token, id)
}

func (e *eventStore) GetPreferences(ctx context.Context, token string, page users.Page) (users.PreferencesPage, error) {
	return e.svc.GetPreferences(ctx, token, page)
}

func (e *eventStore) UpdatePreferences(ctx context.Context, token string, preference users.Preference) (users.Preference, error) {
	return e.svc.UpdatePreferences(ctx, token, preference)
}

func (e *eventStore) UpdateEmailPreferences(ctx context.Context, token string, preference users.Preference) (users.Preference, error) {
	return e.svc.UpdateEmailPreferences(ctx, token, preference)
}

func (e *eventStore) UpdatePushPreferences(ctx context.Context, token string, preference users.Preference) (users.Preference, error) {
	return e.svc.UpdatePushPreferences(ctx, token, preference)
}

func (e *eventStore) DeletePreferences(ctx context.Context, token string) error {
	return e.svc.DeletePreferences(ctx, token)
}

func (e *eventStore) CreateFollower(ctx context.Context, token string, following users.Following) (users.Following, error) {
	following, err := e.svc.CreateFollower(ctx, token, following)
	if err != nil {
		return users.Following{}, err
	}

	if err := e.Publish(ctx, e.publisher, "followers.created", following); err != nil {
		return following, err
	}

	return following, nil
}

func (e *eventStore) GetUserFollowings(ctx context.Context, token string, page users.Page) (users.FollowingsPage, error) {
	return e.svc.GetUserFollowings(ctx, token, page)
}

func (e *eventStore) DeleteFollower(ctx context.Context, token string, following users.Following) error {
	return e.svc.DeleteFollower(ctx, token, following)
}

func (e *eventStore) CreateFeed(ctx context.Context, feed users.Feed) error {
	return e.svc.CreateFeed(ctx, feed)
}

func (e *eventStore) GetUserFeed(ctx context.Context, token string, page users.Page) (users.FeedPage, error) {
	return e.svc.GetUserFeed(ctx, token, page)
}

func (e *eventStore) Publish(ctx context.Context, pub events.Publisher, topic string, entity users.Following) error {
	jsonEntity, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(jsonEntity, &payload); err != nil {
		return err
	}

	return pub.Publish(ctx, topic, payload)
}
