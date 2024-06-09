// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package repository_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/0x6flab/namegenerator"
	"github.com/gofrs/uuid"
	"github.com/rodneyosodo/twiga/notifications"
	"github.com/rodneyosodo/twiga/notifications/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	invalidID   = uuid.Must(uuid.NewV4()).String()
	malformedID = strings.Repeat("a", 100)
	namegen     = namegenerator.NewGenerator()
)

func TestCreateNotification(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM notifications")
		require.NoError(t, err)
	})
	repo := repository.NewRepository(db)

	cases := []struct {
		desc         string
		notification notifications.Notification
		err          error
	}{
		{
			desc: "valid notification",
			notification: notifications.Notification{
				UserID:   uuid.Must(uuid.NewV4()).String(),
				Category: notifications.Post,
				Content:  namegen.Generate(),
			},
			err: nil,
		},
		{
			desc: "empty user id",
			notification: notifications.Notification{
				UserID:   "",
				Category: notifications.Post,
				Content:  namegen.Generate(),
			},
			err: errors.New("invalid input syntax for type uuid"),
		},
		{
			desc: "empty content",
			notification: notifications.Notification{
				UserID:   uuid.Must(uuid.NewV4()).String(),
				Category: notifications.Post,
				Content:  "",
			},
			err: nil,
		},
		{
			desc: "malformed user id",
			notification: notifications.Notification{
				UserID:   malformedID,
				Category: notifications.Post,
				Content:  namegen.Generate(),
			},
			err: errors.New("invalid input syntax for type uuid"),
		},
		{
			desc: "unknown user id",
			notification: notifications.Notification{
				UserID:   invalidID,
				Category: notifications.Post,
				Content:  namegen.Generate(),
			},
			err: nil,
		},
		{
			desc:         "empty notification",
			notification: notifications.Notification{},
			err:          errors.New("invalid input syntax for type uuid"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			saved, err := repo.CreateNotification(context.Background(), tc.notification)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())

				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, saved.ID)
			assert.Equal(t, tc.notification.UserID, saved.UserID)
		})
	}
}

func TestRetrieveNotification(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM notifications")
		require.NoError(t, err)
	})
	repo := repository.NewRepository(db)

	notification := notifications.Notification{
		UserID:   uuid.Must(uuid.NewV4()).String(),
		Category: notifications.Post,
		Content:  namegen.Generate(),
	}
	notification, err := repo.CreateNotification(context.Background(), notification)
	require.NoError(t, err)

	cases := []struct {
		desc     string
		id       string
		response notifications.Notification
		err      error
	}{
		{
			desc:     "valid notification",
			id:       notification.ID,
			response: notification,
			err:      nil,
		},
		{
			desc:     "empty id",
			id:       "",
			response: notifications.Notification{},
			err:      errors.New("invalid input syntax for type uuid"),
		},
		{
			desc:     "malformed id",
			id:       malformedID,
			response: notifications.Notification{},
			err:      errors.New("invalid input syntax for type uuid"),
		},
		{
			desc:     "empty notification",
			id:       uuid.Must(uuid.NewV4()).String(),
			response: notifications.Notification{},
			err:      errors.New("notification not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			n, err := repo.RetrieveNotification(context.Background(), tc.id)
			switch {
			case tc.err != nil:
				assert.ErrorContains(t, err, tc.err.Error())
			default:
				assert.Equal(t, tc.response, n)
			}
		})
	}
}

func TestRetrieveAllNotifications(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM notifications")
		require.NoError(t, err)
	})
	repo := repository.NewRepository(db)

	num := uint64(10)
	saved := make([]notifications.Notification, num)
	for i := range num {
		n := notifications.Notification{
			UserID:   uuid.Must(uuid.NewV4()).String(),
			Category: notifications.Post,
			Content:  namegen.Generate(),
		}
		if i == 0 {
			n.IsRead = true
			n.Category = notifications.Like
		}

		n, err := repo.CreateNotification(context.Background(), n)
		require.NoError(t, err)
		saved[i] = n
	}

	cases := []struct {
		desc     string
		page     notifications.Page
		response notifications.NotificationsPage
		err      error
	}{
		{
			desc: "valid page",
			page: notifications.Page{
				Offset: 0,
				Limit:  10,
			},
			response: notifications.NotificationsPage{
				Page: notifications.Page{
					Total:  num,
					Offset: 0,
					Limit:  10,
				},
				Notifications: saved,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: notifications.Page{
				Offset: num + 1,
				Limit:  10,
			},
			response: notifications.NotificationsPage{
				Page: notifications.Page{
					Total:  num,
					Offset: num + 1,
					Limit:  10,
				},
				Notifications: []notifications.Notification{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: notifications.Page{
				Offset: 0,
				Limit:  num + 1,
			},
			response: notifications.NotificationsPage{
				Page: notifications.Page{
					Total:  num,
					Offset: 0,
					Limit:  num + 1,
				},
				Notifications: saved,
			},
			err: nil,
		},
		{
			desc: "for userID",
			page: notifications.Page{
				Offset: 0,
				Limit:  10,
				UserID: saved[0].UserID,
			},
			response: notifications.NotificationsPage{
				Page: notifications.Page{
					Total:  1,
					Offset: 0,
					Limit:  10,
				},
				Notifications: []notifications.Notification{saved[0]},
			},
			err: nil,
		},
		{
			desc: "for category",
			page: notifications.Page{
				Offset:   0,
				Limit:    10,
				Category: saved[0].Category,
			},
			response: notifications.NotificationsPage{
				Page: notifications.Page{
					Total:  1,
					Offset: 0,
					Limit:  10,
				},
				Notifications: []notifications.Notification{saved[0]},
			},
			err: nil,
		},
		{
			desc: "for is read",
			page: notifications.Page{
				Offset: 0,
				Limit:  10,
				IsRead: &saved[0].IsRead,
			},
			response: notifications.NotificationsPage{
				Page: notifications.Page{
					Total:  1,
					Offset: 0,
					Limit:  10,
				},
				Notifications: []notifications.Notification{saved[0]},
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			n, err := repo.RetrieveAllNotifications(context.Background(), tc.page)
			switch {
			case tc.err != nil:
				assert.ErrorContains(t, err, tc.err.Error())
			default:
				assert.Equal(t, tc.response.Total, n.Total)
				assert.Equal(t, tc.response.Offset, n.Offset)
				assert.Equal(t, tc.response.Limit, n.Limit)
				assert.ElementsMatch(t, tc.response.Notifications, n.Notifications)
			}
		})
	}
}

func TestReadNotification(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM notifications")
		require.NoError(t, err)
	})
	repo := repository.NewRepository(db)

	notification := notifications.Notification{
		UserID:   uuid.Must(uuid.NewV4()).String(),
		Category: notifications.Post,
		Content:  namegen.Generate(),
	}
	notification, err := repo.CreateNotification(context.Background(), notification)
	require.NoError(t, err)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "valid notification",
			id:   notification.ID,
			err:  nil,
		},
		{
			desc: "empty id",
			id:   "",
			err:  errors.New("invalid input syntax for type uuid"),
		},
		{
			desc: "malformed id",
			id:   malformedID,
			err:  errors.New("invalid input syntax for type uuid"),
		},
		{
			desc: "empty notification",
			id:   uuid.Must(uuid.NewV4()).String(),
			err:  errors.New("could not read notification"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.ReadNotification(context.Background(), notification.UserID, tc.id)
			switch {
			case tc.err != nil:
				assert.ErrorContains(t, err, tc.err.Error())
			default:
				n, err := repo.RetrieveNotification(context.Background(), tc.id)
				require.NoError(t, err)
				assert.True(t, n.IsRead)
			}
		})
	}
}

func TestReadAllNotifications(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM notifications")
		require.NoError(t, err)
	})
	repo := repository.NewRepository(db)

	num := uint64(10)
	saved := make([]notifications.Notification, num)
	for i := range num {
		n := notifications.Notification{
			UserID:   uuid.Must(uuid.NewV4()).String(),
			Category: notifications.Post,
			Content:  namegen.Generate(),
		}

		n, err := repo.CreateNotification(context.Background(), n)
		require.NoError(t, err)
		saved[i] = n
	}

	cases := []struct {
		desc string
		page notifications.Page
		err  error
	}{
		{
			desc: "valid page",
			page: notifications.Page{
				Offset: 0,
				Limit:  10,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: notifications.Page{
				Offset: num + 1,
				Limit:  10,
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: notifications.Page{
				Offset: 0,
				Limit:  num + 1,
			},
			err: nil,
		},
		{
			desc: "for userID",
			page: notifications.Page{
				Offset: 0,
				Limit:  10,
				UserID: saved[0].UserID,
			},
			err: nil,
		},
		{
			desc: "for category",
			page: notifications.Page{
				Offset:   0,
				Limit:    10,
				Category: saved[0].Category,
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.ReadAllNotifications(context.Background(), tc.page)
			switch {
			case tc.err != nil:
				assert.ErrorContains(t, err, tc.err.Error())
			default:
				n, err := repo.RetrieveAllNotifications(context.Background(), tc.page)
				require.NoError(t, err)
				for _, n := range n.Notifications {
					assert.True(t, n.IsRead)
				}
			}
		})
	}
}

func TestDeleteNotification(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM notifications")
		require.NoError(t, err)
	})
	repo := repository.NewRepository(db)

	notification := notifications.Notification{
		UserID:   uuid.Must(uuid.NewV4()).String(),
		Category: notifications.Post,
		Content:  namegen.Generate(),
	}
	notification, err := repo.CreateNotification(context.Background(), notification)
	require.NoError(t, err)

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "valid notification",
			id:   notification.ID,
			err:  nil,
		},
		{
			desc: "empty id",
			id:   "",
			err:  errors.New("invalid input syntax for type uuid"),
		},
		{
			desc: "malformed id",
			id:   malformedID,
			err:  errors.New("invalid input syntax for type uuid"),
		},
		{
			desc: "empty notification",
			id:   uuid.Must(uuid.NewV4()).String(),
			err:  errors.New("could not delete notification"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.DeleteNotification(context.Background(), tc.id)
			switch {
			case tc.err != nil:
				assert.ErrorContains(t, err, tc.err.Error())
			default:
				n, err := repo.RetrieveNotification(context.Background(), tc.id)
				assert.ErrorContains(t, err, "notification not found")
				assert.Empty(t, n)
			}
		})
	}
}
