// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/rodneyosodo/twiga/users"
)

var _ users.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger

	svc users.Service
}

func NewLoggingMiddleware(logger *slog.Logger, svc users.Service) users.Service {
	return &loggingMiddleware{
		logger: logger,
		svc:    svc,
	}
}

func (lm *loggingMiddleware) IssueToken(ctx context.Context, user users.User) (token string, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Issue token failed", args...)

			return
		}
		lm.logger.Info("Issue token completed successfully", args...)
	}(time.Now())

	return lm.svc.IssueToken(ctx, user)
}

func (lm *loggingMiddleware) RefreshToken(ctx context.Context, token string) (t string, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Refresh token failed", args...)

			return
		}
		lm.logger.Info("Refresh token completed successfully", args...)
	}(time.Now())

	return lm.svc.RefreshToken(ctx, token)
}

func (lm *loggingMiddleware) IdentifyUser(ctx context.Context, token string) (userID string, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("user_id", userID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Identify user failed", args...)

			return
		}
		lm.logger.Info("Identify user completed successfully", args...)
	}(time.Now())

	return lm.svc.IdentifyUser(ctx, token)
}

func (lm *loggingMiddleware) CreateUser(ctx context.Context, user users.User) (c users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", c.ID),
				slog.String("username", c.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create user failed", args...)

			return
		}
		lm.logger.Info("Create user completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateUser(ctx, user)
}

func (lm *loggingMiddleware) GetUserByID(ctx context.Context, token, id string) (u users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", u.ID),
				slog.String("username", u.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Get user by ID failed", args...)

			return
		}
		lm.logger.Info("Get user by ID completed successfully", args...)
	}(time.Now())

	return lm.svc.GetUserByID(ctx, token, id)
}

func (lm *loggingMiddleware) GetUsers(ctx context.Context, token string, page users.Page) (up users.UsersPage, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("page",
				slog.Uint64("limit", up.Limit),
				slog.Uint64("offset", up.Offset),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Get users failed", args...)

			return
		}
		lm.logger.Info("Get users completed successfully", args...)
	}(time.Now())

	return lm.svc.GetUsers(ctx, token, page)
}

func (lm *loggingMiddleware) UpdateUser(ctx context.Context, token string, user users.User) (u users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", u.ID),
				slog.String("username", u.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update user failed", args...)

			return
		}
		lm.logger.Info("Update user completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateUser(ctx, token, user)
}

func (lm *loggingMiddleware) UpdateUserUsername(ctx context.Context, token string, user users.User) (u users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", u.ID),
				slog.String("username", u.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update user username failed", args...)

			return
		}
		lm.logger.Info("Update user username completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateUserUsername(ctx, token, user)
}

func (lm *loggingMiddleware) UpdateUserPassword(ctx context.Context, token, oldPassword, newPassword string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update user password failed", args...)

			return
		}
		lm.logger.Info("Update user password completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateUserPassword(ctx, token, oldPassword, newPassword)
}

func (lm *loggingMiddleware) UpdateUserEmail(ctx context.Context, token string, user users.User) (u users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", u.ID),
				slog.String("username", u.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update user email failed", args...)

			return
		}
		lm.logger.Info("Update user email completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateUserEmail(ctx, token, user)
}

func (lm *loggingMiddleware) UpdateUserBio(ctx context.Context, token string, user users.User) (u users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", u.ID),
				slog.String("username", u.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update user bio failed", args...)

			return
		}
		lm.logger.Info("Update user bio completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateUserBio(ctx, token, user)
}

func (lm *loggingMiddleware) UpdateUserPictureURL(ctx context.Context, token string, user users.User) (u users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", u.ID),
				slog.String("username", u.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update user picture URL failed", args...)

			return
		}
		lm.logger.Info("Update user picture URL completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateUserPictureURL(ctx, token, user)
}

func (lm *loggingMiddleware) UpdateUserPreferences(ctx context.Context, token string, user users.User) (u users.User, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("user",
				slog.String("id", u.ID),
				slog.String("username", u.Username),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update user preferences failed", args...)

			return
		}
		lm.logger.Info("Update user preferences completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateUserPreferences(ctx, token, user)
}

func (lm *loggingMiddleware) DeleteUser(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("user_id", id),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Delete user failed", args...)

			return
		}
		lm.logger.Info("Delete user completed successfully", args...)
	}(time.Now())

	return lm.svc.DeleteUser(ctx, token, id)
}

func (lm *loggingMiddleware) CreatePreferences(ctx context.Context, token string, preference users.Preference) (p users.Preference, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("preference",
				slog.String("id", p.ID),
				slog.String("user_id", p.UserID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create preferences failed", args...)

			return
		}
		lm.logger.Info("Create preferences completed successfully", args...)
	}(time.Now())

	return lm.svc.CreatePreferences(ctx, token, preference)
}

func (lm *loggingMiddleware) GetPreferencesByUserID(ctx context.Context, token, id string) (p users.Preference, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("user_id", p.UserID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Get preferences by user ID failed", args...)

			return
		}
		lm.logger.Info("Get preferences by user ID completed successfully", args...)
	}(time.Now())

	return lm.svc.GetPreferencesByUserID(ctx, token, id)
}

func (lm *loggingMiddleware) GetPreferences(ctx context.Context, token string, page users.Page) (pp users.PreferencesPage, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("page",
				slog.Uint64("limit", pp.Limit),
				slog.Uint64("offset", pp.Offset),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Get preferences failed", args...)

			return
		}
		lm.logger.Info("Get preferences completed successfully", args...)
	}(time.Now())

	return lm.svc.GetPreferences(ctx, token, page)
}

func (lm *loggingMiddleware) UpdatePreferences(ctx context.Context, token string, preference users.Preference) (p users.Preference, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("preference",
				slog.String("id", p.ID),
				slog.String("user_id", p.UserID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update preferences failed", args...)

			return
		}
		lm.logger.Info("Update preferences completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdatePreferences(ctx, token, preference)
}

func (lm *loggingMiddleware) UpdateEmailPreferences(ctx context.Context, token string, preference users.Preference) (p users.Preference, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("preference",
				slog.String("id", p.ID),
				slog.String("user_id", p.UserID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update email preferences failed", args...)

			return
		}
		lm.logger.Info("Update email preferences completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateEmailPreferences(ctx, token, preference)
}

func (lm *loggingMiddleware) UpdatePushPreferences(ctx context.Context, token string, preference users.Preference) (p users.Preference, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("preference",
				slog.String("id", p.ID),
				slog.String("user_id", p.UserID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update push preferences failed", args...)

			return
		}
		lm.logger.Info("Update push preferences completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdatePushPreferences(ctx, token, preference)
}

func (lm *loggingMiddleware) DeletePreferences(ctx context.Context, token string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Delete preferences failed", args...)

			return
		}
		lm.logger.Info("Delete preferences completed successfully", args...)
	}(time.Now())

	return lm.svc.DeletePreferences(ctx, token)
}

func (lm *loggingMiddleware) CreateFollower(ctx context.Context, token string, following users.Following) (f users.Following, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("following",
				slog.String("id", f.ID),
				slog.String("follower_id", f.FollowerID),
				slog.String("followee_id", f.FolloweeID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create follower failed", args...)

			return
		}
		lm.logger.Info("Create follower completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateFollower(ctx, token, following)
}

func (lm *loggingMiddleware) GetUserFollowings(ctx context.Context, token string, page users.Page) (fp users.FollowingsPage, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("page",
				slog.Uint64("limit", fp.Limit),
				slog.Uint64("offset", fp.Offset),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Get user followings failed", args...)

			return
		}
		lm.logger.Info("Get user followings completed successfully", args...)
	}(time.Now())

	return lm.svc.GetUserFollowings(ctx, token, page)
}

func (lm *loggingMiddleware) DeleteFollower(ctx context.Context, token string, following users.Following) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("following",
				slog.String("id", following.ID),
				slog.String("follower_id", following.FollowerID),
				slog.String("followee_id", following.FolloweeID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Delete follower failed", args...)

			return
		}
		lm.logger.Info("Delete follower completed successfully", args...)
	}(time.Now())

	return lm.svc.DeleteFollower(ctx, token, following)
}

func (lm *loggingMiddleware) CreateFeed(ctx context.Context, feed users.Feed) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("feed",
				slog.String("id", feed.ID),
				slog.String("user_id", feed.UserID),
				slog.String("post_id", feed.PostID),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create feed failed", args...)

			return
		}
		lm.logger.Info("Create feed completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateFeed(ctx, feed)
}

func (lm *loggingMiddleware) GetUserFeed(ctx context.Context, token string, page users.Page) (fp users.FeedPage, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Group("page",
				slog.Uint64("limit", fp.Limit),
				slog.Uint64("offset", fp.Offset),
			),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Get user feed failed", args...)

			return
		}
		lm.logger.Info("Get user feed completed successfully", args...)
	}(time.Now())

	return lm.svc.GetUserFeed(ctx, token, page)
}
