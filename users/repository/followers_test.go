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

	"github.com/rodneyosodo/twiga/users"
	"github.com/rodneyosodo/twiga/users/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFollowers(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM followers")
		require.NoError(t, err)
	})
	repo := repository.NewFollowingRepository(db)
	urepo := repository.NewUsersRepository(db)

	savedUser1, err := urepo.Create(context.Background(), generateUser())
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	savedUser2, err := urepo.Create(context.Background(), generateUser())
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc      string
		following users.Following
		err       error
	}{
		{
			desc: "valid following",
			following: users.Following{
				FollowerID: savedUser1.ID,
				FolloweeID: savedUser2.ID,
			},
			err: nil,
		},
		{
			desc: "same ids",
			following: users.Following{
				FollowerID: savedUser1.ID,
				FolloweeID: savedUser1.ID,
			},
			err: errors.New("new row for relation"),
		},
		{
			desc:      "empty following",
			following: users.Following{},
			err:       errors.New("invalid input syntax"),
		},
		{
			desc: "unknown following",
			following: users.Following{
				FollowerID: invalidID,
				FolloweeID: savedUser2.ID,
			},
			err: errors.New("insert or update on table"),
		},
		{
			desc: "malformed following",
			following: users.Following{
				FollowerID: malformedID,
				FolloweeID: savedUser2.ID,
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			following, err := repo.Create(context.Background(), tc.following)
			switch {
			case err == nil:
				assert.Equal(t, tc.following.FolloweeID, following.FolloweeID)
				assert.Equal(t, tc.following.FolloweeID, following.FolloweeID)
			default:
				assert.Empty(t, following)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestRetrieveAllFollowing(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM followers")
		require.NoError(t, err)
	})
	repo := repository.NewFollowingRepository(db)
	urepo := repository.NewUsersRepository(db)

	num := uint64(10)
	followers := make([]users.User, num)
	followees := make([]users.User, num)
	for i := range num {
		savedUser, err := urepo.Create(context.Background(), generateUser())
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		followers[i] = savedUser
		savedUser, err = urepo.Create(context.Background(), generateUser())
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		followees[i] = savedUser
	}

	followings := make([]users.Following, 0)
	for follower := range followers {
		for followee := range followees {
			following, err := repo.Create(context.Background(), users.Following{
				FollowerID: followers[follower].ID,
				FolloweeID: followees[followee].ID,
			})
			if err != nil {
				t.Fatalf("failed to create following: %v", err)
			}
			followings = append(followings, following)
		}
	}

	cases := []struct {
		desc     string
		page     users.Page
		response users.FollowingsPage
		err      error
	}{
		{
			desc: "valid page",
			page: users.Page{
				Offset: 0,
				Limit:  10,
			},
			response: users.FollowingsPage{
				Page: users.Page{
					Offset: 0,
					Limit:  10,
					Total:  uint64(len(followings)),
				},
				Followings: followings,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: users.Page{
				Offset: num + 1,
				Limit:  10,
			},
			response: users.FollowingsPage{
				Page: users.Page{
					Offset: num + 1,
					Limit:  10,
					Total:  uint64(len(followings)),
				},
				Followings: []users.Following{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: users.Page{
				Offset: 0,
				Limit:  num + 1,
			},
			response: users.FollowingsPage{
				Page: users.Page{
					Offset: 0,
					Limit:  num + 1,
					Total:  uint64(len(followings)),
				},
				Followings: followings,
			},
			err: nil,
		},
		{
			desc: "follower filter",
			page: users.Page{
				Offset:     0,
				Limit:      10,
				FollowerID: followers[0].ID,
			},
			response: users.FollowingsPage{
				Page: users.Page{
					Offset: 0,
					Limit:  10,
					Total:  count(followings, followers[0].ID, true),
				},
				Followings: follows(followings, followers[0].ID, true),
			},
			err: nil,
		},
		{
			desc: "followee filter",
			page: users.Page{
				Offset:     0,
				Limit:      10,
				FolloweeID: followers[0].ID,
			},
			response: users.FollowingsPage{
				Page: users.Page{
					Offset: 0,
					Limit:  10,
					Total:  count(followings, followers[0].ID, false),
				},
				Followings: follows(followings, followers[0].ID, false),
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

func TestDeleteFollowers(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM followers")
		require.NoError(t, err)
	})
	repo := repository.NewFollowingRepository(db)
	urepo := repository.NewUsersRepository(db)

	savedUser1, err := urepo.Create(context.Background(), generateUser())
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	savedUser2, err := urepo.Create(context.Background(), generateUser())
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	following, err := repo.Create(context.Background(), users.Following{
		FollowerID: savedUser1.ID,
		FolloweeID: savedUser2.ID,
	})
	if err != nil {
		t.Fatalf("failed to create following: %v", err)
	}

	cases := []struct {
		desc      string
		following users.Following
		err       error
	}{
		{
			desc:      "valid following",
			following: following,
			err:       nil,
		},
		{
			desc: "invalid following",
			following: users.Following{
				FollowerID: invalidID,
				FolloweeID: savedUser2.ID,
			},
			err: errors.New("following not found"),
		},
		{
			desc: "malformed following",
			following: users.Following{
				FollowerID: malformedID,
				FolloweeID: savedUser2.ID,
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.Delete(context.Background(), tc.following)
			if err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func count(followings []users.Following, id string, follow bool) uint64 {
	var following uint64
	var followers uint64
	for _, f := range followings {
		if f.FollowerID == id {
			followers++
		}
		if f.FolloweeID == id {
			following++
		}
	}

	if follow {
		return following
	}

	return followers
}

func follows(followings []users.Following, id string, follow bool) []users.Following {
	var following uint64
	var followers uint64
	for _, f := range followings {
		if f.FollowerID == id {
			followers++
		}
		if f.FolloweeID == id {
			following++
		}
	}

	if follow {
		return followings[:following]
	}

	return followings[:followers]
}
