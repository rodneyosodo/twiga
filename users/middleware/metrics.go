// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/rodneyosodo/twiga/users"
)

var _ users.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     users.Service
}

func NewMetricsMiddleware(counter metrics.Counter, latency metrics.Histogram, svc users.Service) users.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (m *metricsMiddleware) IssueToken(ctx context.Context, user users.User) (string, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "issue_token").Add(1)
		m.latency.With("method", "issue_token").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.IssueToken(ctx, user)
}

func (m *metricsMiddleware) RefreshToken(ctx context.Context, token string) (string, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "refresh_token").Add(1)
		m.latency.With("method", "refresh_token").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.RefreshToken(ctx, token)
}

func (m *metricsMiddleware) IdentifyUser(ctx context.Context, token string) (userID string, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "identify_user").Add(1)
		m.latency.With("method", "identify_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.IdentifyUser(ctx, token)
}

func (m *metricsMiddleware) CreateUser(ctx context.Context, user users.User) (saved users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "create_user").Add(1)
		m.latency.With("method", "create_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.CreateUser(ctx, user)
}

func (m *metricsMiddleware) GetUserByID(ctx context.Context, token, id string) (user users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "get_user_by_id").Add(1)
		m.latency.With("method", "get_user_by_id").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.GetUserByID(ctx, token, id)
}

func (m *metricsMiddleware) GetUsers(ctx context.Context, token string, page users.Page) (users.UsersPage, error) {
	defer func(begin time.Time) {
		m.counter.With("method", "get_users").Add(1)
		m.latency.With("method", "get_users").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.GetUsers(ctx, token, page)
}

func (m *metricsMiddleware) UpdateUser(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_user").Add(1)
		m.latency.With("method", "update_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateUser(ctx, token, user)
}

func (m *metricsMiddleware) UpdateUserUsername(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_user_username").Add(1)
		m.latency.With("method", "update_user_username").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateUserUsername(ctx, token, user)
}

func (m *metricsMiddleware) UpdateUserPassword(ctx context.Context, token string, oldPassword, currentPassowrd string) (updated users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_user_password").Add(1)
		m.latency.With("method", "update_user_password").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateUserPassword(ctx, token, oldPassword, currentPassowrd)
}

func (m *metricsMiddleware) UpdateUserEmail(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_user_email").Add(1)
		m.latency.With("method", "update_user_email").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateUserEmail(ctx, token, user)
}

func (m *metricsMiddleware) UpdateUserBio(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_user_bio").Add(1)
		m.latency.With("method", "update_user_bio").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateUserBio(ctx, token, user)
}

func (m *metricsMiddleware) UpdateUserPictureURL(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_user_picture_url").Add(1)
		m.latency.With("method", "update_user_picture_url").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateUserPictureURL(ctx, token, user)
}

func (m *metricsMiddleware) UpdateUserPreferences(ctx context.Context, token string, user users.User) (updated users.User, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_user_preferences").Add(1)
		m.latency.With("method", "update_user_preferences").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateUserPreferences(ctx, token, user)
}

func (m *metricsMiddleware) DeleteUser(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		m.counter.With("method", "delete_user").Add(1)
		m.latency.With("method", "delete_user").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.DeleteUser(ctx, token, id)
}

func (m *metricsMiddleware) CreatePreferences(ctx context.Context, token string, preference users.Preference) (created users.Preference, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "create_preferences").Add(1)
		m.latency.With("method", "create_preferences").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.CreatePreferences(ctx, token, preference)
}

func (m *metricsMiddleware) GetPreferencesByUserID(ctx context.Context, token string) (preference users.Preference, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "get_preferences_by_user_id").Add(1)
		m.latency.With("method", "get_preferences_by_user_id").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.GetPreferencesByUserID(ctx, token)
}

func (m *metricsMiddleware) GetPreferences(ctx context.Context, token string, page users.Page) (preferences users.PreferencesPage, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "get_preferences").Add(1)
		m.latency.With("method", "get_preferences").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.GetPreferences(ctx, token, page)
}

func (m *metricsMiddleware) UpdatePreferences(ctx context.Context, token string, preference users.Preference) (updated users.Preference, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_preferences").Add(1)
		m.latency.With("method", "update_preferences").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdatePreferences(ctx, token, preference)
}

func (m *metricsMiddleware) UpdateEmailPreferences(ctx context.Context, token string, preference users.Preference) (updated users.Preference, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_email_preferences").Add(1)
		m.latency.With("method", "update_email_preferences").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdateEmailPreferences(ctx, token, preference)
}

func (m *metricsMiddleware) UpdatePushPreferences(ctx context.Context, token string, preference users.Preference) (updated users.Preference, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "update_push_preferences").Add(1)
		m.latency.With("method", "update_push_preferences").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.UpdatePushPreferences(ctx, token, preference)
}

func (m *metricsMiddleware) DeletePreferences(ctx context.Context, token string) error {
	defer func(begin time.Time) {
		m.counter.With("method", "delete_preferences").Add(1)
		m.latency.With("method", "delete_preferences").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.DeletePreferences(ctx, token)
}

func (m *metricsMiddleware) CreateFollower(ctx context.Context, token string, following users.Following) (created users.Following, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "create_follower").Add(1)
		m.latency.With("method", "create_follower").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.CreateFollower(ctx, token, following)
}

func (m *metricsMiddleware) GetUserFollowings(ctx context.Context, token string, page users.Page) (followings users.FollowingsPage, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "get_user_followings").Add(1)
		m.latency.With("method", "get_user_followings").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.GetUserFollowings(ctx, token, page)
}

func (m *metricsMiddleware) DeleteFollower(ctx context.Context, token string, following users.Following) error {
	defer func(begin time.Time) {
		m.counter.With("method", "delete_follower").Add(1)
		m.latency.With("method", "delete_follower").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.DeleteFollower(ctx, token, following)
}

func (m *metricsMiddleware) CreateFeed(ctx context.Context, feed users.Feed) error {
	defer func(begin time.Time) {
		m.counter.With("method", "create_feed").Add(1)
		m.latency.With("method", "create_feed").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.CreateFeed(ctx, feed)
}

func (m *metricsMiddleware) GetUserFeed(ctx context.Context, token string, page users.Page) (feeds users.FeedPage, err error) {
	defer func(begin time.Time) {
		m.counter.With("method", "get_user_feed").Add(1)
		m.latency.With("method", "get_user_feed").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.svc.GetUserFeed(ctx, token, page)
}
