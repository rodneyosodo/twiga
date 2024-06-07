// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package middleware

import (
	"context"

	"github.com/rodneyosodo/twiga/users"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var _ users.Service = (*tracingMiddleware)(nil)

type tracingMiddleware struct {
	tracer trace.Tracer
	svc    users.Service
}

func NewTracingMiddleware(tracer trace.Tracer, svc users.Service) users.Service {
	return &tracingMiddleware{
		tracer: tracer,
		svc:    svc,
	}
}

func (m *tracingMiddleware) IssueToken(ctx context.Context, user users.User) (string, error) {
	ctx, span := m.tracer.Start(ctx, "issue_token")
	defer span.End()

	return m.svc.IssueToken(ctx, user)
}

func (m *tracingMiddleware) RefreshToken(ctx context.Context, token string) (string, error) {
	ctx, span := m.tracer.Start(ctx, "refresh_token")
	defer span.End()

	return m.svc.RefreshToken(ctx, token)
}

func (m *tracingMiddleware) IdentifyUser(ctx context.Context, token string) (userID string, err error) {
	ctx, span := m.tracer.Start(ctx, "identify_user", trace.WithAttributes(
		attribute.String("user_id", userID),
	))
	defer span.End()

	return m.svc.IdentifyUser(ctx, token)
}

func (m *tracingMiddleware) CreateUser(ctx context.Context, user users.User) (saved users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "create_user", trace.WithAttributes(
		attribute.String("user_id", saved.ID),
		attribute.String("user_name", saved.Username),
		attribute.String("user_display_name", saved.DisplayName),
	))
	defer span.End()

	return m.svc.CreateUser(ctx, user)
}

func (m *tracingMiddleware) GetUserByID(ctx context.Context, token, id string) (user users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "get_user_by_id", trace.WithAttributes(
		attribute.String("user_id", id),
	))
	defer span.End()

	return m.svc.GetUserByID(ctx, token, id)
}

func (m *tracingMiddleware) GetUsers(ctx context.Context, token string, page users.Page) (users.UsersPage, error) {
	ctx, span := m.tracer.Start(ctx, "get_users")
	defer span.End()

	return m.svc.GetUsers(ctx, token, page)
}

func (m *tracingMiddleware) UpdateUser(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "update_user", trace.WithAttributes(
		attribute.String("user_id", user.ID),
		attribute.String("user_name", user.Username),
		attribute.String("user_display_name", user.DisplayName),
	))
	defer span.End()

	return m.svc.UpdateUser(ctx, token, user)
}

func (m *tracingMiddleware) UpdateUserUsername(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "update_user_username", trace.WithAttributes(
		attribute.String("user_id", user.ID),
		attribute.String("user_name", user.Username),
	))
	defer span.End()

	return m.svc.UpdateUserUsername(ctx, token, user)
}

func (m *tracingMiddleware) UpdateUserPassword(ctx context.Context, token, oldPassword, newPassword string) (updated users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "update_user_password", trace.WithAttributes(
		attribute.String("user_id", updated.ID),
	))
	defer span.End()

	return m.svc.UpdateUserPassword(ctx, token, oldPassword, newPassword)
}

func (m *tracingMiddleware) UpdateUserEmail(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "update_user_email", trace.WithAttributes(
		attribute.String("user_id", user.ID),
	))
	defer span.End()

	return m.svc.UpdateUserEmail(ctx, token, user)
}

func (m *tracingMiddleware) UpdateUserBio(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "update_user_bio", trace.WithAttributes(
		attribute.String("user_id", user.ID),
		attribute.String("user_bio", user.Bio),
	))
	defer span.End()

	return m.svc.UpdateUserBio(ctx, token, user)
}

func (m *tracingMiddleware) UpdateUserPictureURL(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "update_user_picture_url", trace.WithAttributes(
		attribute.String("user_id", user.ID),
		attribute.String("user_picture_url", user.PictureURL),
	))
	defer span.End()

	return m.svc.UpdateUserPictureURL(ctx, token, user)
}

func (m *tracingMiddleware) UpdateUserPreferences(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	ctx, span := m.tracer.Start(ctx, "update_user_preferences", trace.WithAttributes(
		attribute.String("user_id", user.ID),
		attribute.StringSlice("user_preferences", user.Preferences),
	))
	defer span.End()

	return m.svc.UpdateUserPreferences(ctx, token, user)
}

func (m *tracingMiddleware) DeleteUser(ctx context.Context, token, id string) (err error) {
	ctx, span := m.tracer.Start(ctx, "delete_user", trace.WithAttributes(
		attribute.String("user_id", id),
	))
	defer span.End()

	return m.svc.DeleteUser(ctx, token, id)
}

func (m *tracingMiddleware) CreatePreferences(ctx context.Context, token string, preference users.Preference) (saved users.Preference, err error) {
	ctx, span := m.tracer.Start(ctx, "create_preferences", trace.WithAttributes(
		attribute.String("user_id", preference.UserID),
		attribute.Bool("email_enabled", preference.EmailEnable),
		attribute.Bool("push_enabled", preference.PushEnable),
	))
	defer span.End()

	return m.svc.CreatePreferences(ctx, token, preference)
}

func (m *tracingMiddleware) GetPreferencesByUserID(ctx context.Context, token string) (preference users.Preference, err error) {
	ctx, span := m.tracer.Start(ctx, "get_preferences_by_user_id")
	defer span.End()

	return m.svc.GetPreferencesByUserID(ctx, token)
}

func (m *tracingMiddleware) GetPreferences(ctx context.Context, token string, page users.Page) (preferences users.PreferencesPage, err error) {
	ctx, span := m.tracer.Start(ctx, "get_preferences")
	defer span.End()

	return m.svc.GetPreferences(ctx, token, page)
}

func (m *tracingMiddleware) UpdatePreferences(ctx context.Context, token string, preference users.Preference) (updated users.Preference, err error) {
	ctx, span := m.tracer.Start(ctx, "update_preferences", trace.WithAttributes(
		attribute.String("user_id", preference.UserID),
		attribute.Bool("email_enabled", preference.EmailEnable),
		attribute.Bool("push_enabled", preference.PushEnable),
	))
	defer span.End()

	return m.svc.UpdatePreferences(ctx, token, preference)
}

func (m *tracingMiddleware) UpdateEmailPreferences(ctx context.Context, token string, preference users.Preference) (updated users.Preference, err error) {
	ctx, span := m.tracer.Start(ctx, "update_email_preferences", trace.WithAttributes(
		attribute.String("user_id", preference.UserID),
		attribute.Bool("email_enabled", preference.EmailEnable),
	))
	defer span.End()

	return m.svc.UpdateEmailPreferences(ctx, token, preference)
}

func (m *tracingMiddleware) UpdatePushPreferences(ctx context.Context, token string, preference users.Preference) (updated users.Preference, err error) {
	ctx, span := m.tracer.Start(ctx, "update_push_preferences", trace.WithAttributes(
		attribute.String("user_id", preference.UserID),
		attribute.Bool("push_enabled", preference.PushEnable),
	))
	defer span.End()

	return m.svc.UpdatePushPreferences(ctx, token, preference)
}

func (m *tracingMiddleware) DeletePreferences(ctx context.Context, token string) (err error) {
	ctx, span := m.tracer.Start(ctx, "delete_preferences")
	defer span.End()

	return m.svc.DeletePreferences(ctx, token)
}

func (m *tracingMiddleware) CreateFollower(ctx context.Context, token string, following users.Following) (saved users.Following, err error) {
	ctx, span := m.tracer.Start(ctx, "create_follower", trace.WithAttributes(
		attribute.String("follower_id", following.FollowerID),
		attribute.String("followee_id", following.FolloweeID),
	))
	defer span.End()

	return m.svc.CreateFollower(ctx, token, following)
}

func (m *tracingMiddleware) GetUserFollowings(ctx context.Context, token string, page users.Page) (followings users.FollowingsPage, err error) {
	ctx, span := m.tracer.Start(ctx, "get_user_followings")
	defer span.End()

	return m.svc.GetUserFollowings(ctx, token, page)
}

func (m *tracingMiddleware) DeleteFollower(ctx context.Context, token string, following users.Following) (err error) {
	ctx, span := m.tracer.Start(ctx, "delete_follower", trace.WithAttributes(
		attribute.String("follower_id", following.FollowerID),
		attribute.String("followee_id", following.FolloweeID),
	))
	defer span.End()

	return m.svc.DeleteFollower(ctx, token, following)
}

func (m *tracingMiddleware) CreateFeed(ctx context.Context, feed users.Feed) (err error) {
	ctx, span := m.tracer.Start(ctx, "create_feed", trace.WithAttributes(
		attribute.String("user_id", feed.UserID),
		attribute.String("post_id", feed.PostID),
	))
	defer span.End()

	return m.svc.CreateFeed(ctx, feed)
}

func (m *tracingMiddleware) GetUserFeed(ctx context.Context, token string, page users.Page) (feed users.FeedPage, err error) {
	ctx, span := m.tracer.Start(ctx, "get_user_feed")
	defer span.End()

	return m.svc.GetUserFeed(ctx, token, page)
}
