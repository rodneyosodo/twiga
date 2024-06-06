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
	"testing"

	"github.com/gofrs/uuid"
	"github.com/rodneyosodo/twiga/users"
	"github.com/rodneyosodo/twiga/users/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFeed(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM feeds")
		require.NoError(t, err)
	})
	repo := repository.NewFeedRepository(db)
	urepo := repository.NewUsersRepository(db)

	savedUser, err := urepo.Create(context.Background(), generateUser())
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		feed users.Feed
		err  error
	}{
		{
			desc: "valid feed",
			feed: users.Feed{
				UserID: savedUser.ID,
				PostID: invalidID,
			},
			err: nil,
		},
		{
			desc: "same ids",
			feed: users.Feed{
				UserID: savedUser.ID,
				PostID: savedUser.ID,
			},
			err: errors.New("new row for relation"),
		},
		{
			desc: "empty feed",
			feed: users.Feed{},
			err:  errors.New("invalid input syntax"),
		},
		{
			desc: "unknown feed",
			feed: users.Feed{
				UserID: invalidID,
				PostID: "post_id",
			},
			err: errors.New("invalid input syntax"),
		},
		{
			desc: "malformed feed",
			feed: users.Feed{
				UserID: malformedID,
				PostID: "post_id",
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.Create(context.Background(), tc.feed)
			if err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestRetrieveAllFeed(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM feeds")
		require.NoError(t, err)
	})
	repo := repository.NewFeedRepository(db)
	urepo := repository.NewUsersRepository(db)

	num := uint64(10)
	feeds := make([]users.Feed, num)
	for i := range num {
		savedUser, err := urepo.Create(context.Background(), generateUser())
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		feed := users.Feed{
			UserID: savedUser.ID,
			PostID: uuid.Must(uuid.NewV4()).String(),
		}
		if err := repo.Create(context.Background(), feed); err != nil {
			t.Fatalf("failed to create feed: %v", err)
		}
		feeds[i] = feed
	}

	cases := []struct {
		desc     string
		page     users.Page
		response users.FeedPage
		err      error
	}{
		{
			desc: "valid page",
			page: users.Page{
				Offset: 0,
				Limit:  10,
			},
			response: users.FeedPage{
				Page: users.Page{
					Offset: 0,
					Limit:  10,
					Total:  uint64(len(feeds)),
				},
				Feeds: feeds,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: users.Page{
				Offset: num + 1,
				Limit:  10,
			},
			response: users.FeedPage{
				Page: users.Page{
					Offset: num + 1,
					Limit:  10,
					Total:  uint64(len(feeds)),
				},
				Feeds: []users.Feed{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: users.Page{
				Offset: 0,
				Limit:  num + 1,
			},
			response: users.FeedPage{
				Page: users.Page{
					Offset: 0,
					Limit:  num + 1,
					Total:  uint64(len(feeds)),
				},
				Feeds: feeds,
			},
			err: nil,
		},
		{
			desc: "user id filter",
			page: users.Page{
				Offset: 0,
				Limit:  10,
				UserID: feeds[0].UserID,
			},
			response: users.FeedPage{
				Page: users.Page{
					Offset: 0,
					Limit:  10,
					Total:  countUser(feeds, feeds[0].UserID),
				},
				Feeds: getUserFeeds(feeds, feeds[0].UserID),
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			preferencesPage, err := repo.RetrieveAll(context.Background(), tc.page)
			switch {
			case err == nil:
				assert.Equal(t, tc.response.Offset, preferencesPage.Offset)
				assert.Equal(t, tc.response.Limit, preferencesPage.Limit)
				assert.Equal(t, tc.response.Total, preferencesPage.Total)
			default:
				assert.Empty(t, preferencesPage)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func countUser(feeds []users.Feed, id string) uint64 {
	var count uint64
	for _, feed := range feeds {
		if feed.UserID == id {
			count++
		}
	}

	return count
}

func getUserFeeds(feeds []users.Feed, id string) []users.Feed {
	var items []users.Feed
	for _, feed := range feeds {
		if feed.UserID == id {
			items = append(items, feed)
		}
	}

	return items
}
