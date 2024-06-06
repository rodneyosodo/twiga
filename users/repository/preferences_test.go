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
	"math/rand"
	"strings"
	"testing"

	"github.com/rodneyosodo/twiga/users"
	"github.com/rodneyosodo/twiga/users/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePreference(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM preferences")
		require.NoError(t, err)
	})
	repo := repository.NewPreferencesRepository(db)
	urepo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	savedUser, err := urepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc       string
		preference users.Preference
		err        error
	}{
		{
			desc: "valid preference",
			preference: users.Preference{
				UserID:      savedUser.ID,
				EmailEnable: rand.Intn(2) == 0,
				PushEnable:  rand.Intn(2) == 0,
			},
			err: nil,
		},
		{
			desc:       "empty preference",
			preference: users.Preference{},
			err:        errors.New("invalid input syntax"),
		},
		{
			desc: "unknown user",
			preference: users.Preference{
				UserID:      invalidID,
				EmailEnable: rand.Intn(2) == 0,
				PushEnable:  rand.Intn(2) == 0,
			},
			err: errors.New("insert or update on table"),
		},
		{
			desc: "malformed preference",
			preference: users.Preference{
				UserID:      malformedID,
				EmailEnable: rand.Intn(2) == 0,
				PushEnable:  rand.Intn(2) == 0,
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			preference, err := repo.Create(context.Background(), tc.preference)
			switch {
			case err == nil:
				assert.Equal(t, tc.preference.EmailEnable, preference.EmailEnable)
				assert.Equal(t, tc.preference.PushEnable, preference.PushEnable)
			default:
				assert.Empty(t, preference)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestRetrievePreferenceByUserID(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM preferences")
		require.NoError(t, err)
	})
	repo := repository.NewPreferencesRepository(db)
	urepo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	savedUser, err := urepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	preference := users.Preference{
		UserID:      savedUser.ID,
		EmailEnable: rand.Intn(2) == 0,
		PushEnable:  rand.Intn(2) == 0,
	}
	saved, err := repo.Create(context.Background(), preference)
	if err != nil {
		t.Fatalf("failed to create preference: %v", err)
	}

	cases := []struct {
		desc       string
		id         string
		preference users.Preference
		err        error
	}{
		{
			desc:       "valid preference",
			id:         saved.UserID,
			preference: saved,
			err:        nil,
		},
		{
			desc:       "invalid preference",
			id:         invalidID,
			preference: users.Preference{},
			err:        errors.New("preference not found"),
		},
		{
			desc:       "malformed preference",
			id:         malformedID,
			preference: users.Preference{},
			err:        errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.RetrieveByUserID(context.Background(), tc.id)
			if err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			}
			assert.Equal(t, tc.preference, user)
		})
	}
}

func TestRetrieveAllPreferences(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM preferences")
		require.NoError(t, err)
	})
	repo := repository.NewPreferencesRepository(db)
	urepo := repository.NewUsersRepository(db)

	num := uint64(10)
	savedPreferences := make([]users.Preference, num)

	for i := range num {
		user := users.User{
			Username:    namegen.Generate(),
			DisplayName: namegen.Generate(),
			Bio:         strings.Repeat("a", 100),
			PictureURL:  "https://example.com" + namegen.Generate(),
			Email:       namegen.Generate() + "@example.com",
			Password:    "password",
			Preferences: []string{"test"},
		}
		savedUser, err := urepo.Create(context.Background(), user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		preference := users.Preference{
			UserID:      savedUser.ID,
			EmailEnable: rand.Intn(2) == 0,
			PushEnable:  rand.Intn(2) == 0,
		}
		saved, err := repo.Create(context.Background(), preference)
		if err != nil {
			t.Fatalf("failed to create preference: %v", err)
		}
		savedPreferences[i] = saved
	}

	cases := []struct { //nolint:dupl
		desc     string
		page     users.Page
		response users.PreferencesPage
		err      error
	}{
		{
			desc: "valid page",
			page: users.Page{
				Offset: 0,
				Limit:  10,
			},
			response: users.PreferencesPage{
				Page: users.Page{
					Offset: 0,
					Limit:  10,
					Total:  uint64(len(savedPreferences)),
				},
				Preferences: savedPreferences,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: users.Page{
				Offset: num + 1,
				Limit:  10,
			},
			response: users.PreferencesPage{
				Page: users.Page{
					Offset: num + 1,
					Limit:  10,
					Total:  uint64(len(savedPreferences)),
				},
				Preferences: []users.Preference{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: users.Page{
				Offset: 0,
				Limit:  num + 1,
			},
			response: users.PreferencesPage{
				Page: users.Page{
					Offset: 0,
					Limit:  num + 1,
					Total:  uint64(len(savedPreferences)),
				},
				Preferences: savedPreferences,
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
				assert.ElementsMatch(t, tc.response.Preferences, preferencesPage.Preferences)
			default:
				assert.Empty(t, preferencesPage)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdatePreference(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM preferences")
		require.NoError(t, err)
	})
	repo := repository.NewPreferencesRepository(db)
	urepo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	savedUser, err := urepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	preference := users.Preference{
		UserID:      savedUser.ID,
		EmailEnable: rand.Intn(2) == 0,
		PushEnable:  rand.Intn(2) == 0,
	}
	if _, err := repo.Create(context.Background(), preference); err != nil {
		t.Fatalf("failed to create preference: %v", err)
	}

	cases := []struct {
		desc       string
		preference users.Preference
		err        error
	}{
		{
			desc: "valid preference",
			preference: users.Preference{
				UserID:      savedUser.ID,
				EmailEnable: rand.Intn(2) == 0,
				PushEnable:  rand.Intn(2) == 0,
			},
			err: nil,
		},
		{
			desc: "invalid preference",
			preference: users.Preference{
				UserID:      invalidID,
				EmailEnable: rand.Intn(2) == 0,
				PushEnable:  rand.Intn(2) == 0,
			},
			err: errors.New("preference not found"),
		},
		{
			desc: "malformed preference",
			preference: users.Preference{
				UserID:      malformedID,
				EmailEnable: rand.Intn(2) == 0,
				PushEnable:  rand.Intn(2) == 0,
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			preference, err := repo.Update(context.Background(), tc.preference)
			switch {
			case err == nil:
				assert.Equal(t, tc.preference.EmailEnable, preference.EmailEnable)
				assert.Equal(t, tc.preference.PushEnable, preference.PushEnable)
			default:
				assert.Empty(t, preference)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdatePreferenceEmail(t *testing.T) { //nolint:dupl
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM preferences")
		require.NoError(t, err)
	})
	repo := repository.NewPreferencesRepository(db)
	urepo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	savedUser, err := urepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	preference := users.Preference{
		UserID:      savedUser.ID,
		EmailEnable: rand.Intn(2) == 0,
		PushEnable:  rand.Intn(2) == 0,
	}
	if _, err := repo.Create(context.Background(), preference); err != nil {
		t.Fatalf("failed to create preference: %v", err)
	}

	cases := []struct {
		desc       string
		preference users.Preference
		err        error
	}{
		{
			desc: "valid preference",
			preference: users.Preference{
				UserID:      savedUser.ID,
				EmailEnable: rand.Intn(2) == 0,
			},
			err: nil,
		},
		{
			desc: "invalid preference",
			preference: users.Preference{
				UserID:      invalidID,
				EmailEnable: rand.Intn(2) == 0,
			},
			err: errors.New("preference not found"),
		},
		{
			desc: "malformed preference",
			preference: users.Preference{
				UserID:      malformedID,
				EmailEnable: rand.Intn(2) == 0,
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			preference, err := repo.UpdateEmail(context.Background(), tc.preference)
			switch {
			case err == nil:
				assert.Equal(t, tc.preference.EmailEnable, preference.EmailEnable)
			default:
				assert.Empty(t, preference)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdatePreferencePush(t *testing.T) { //nolint:dupl
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM preferences")
		require.NoError(t, err)
	})
	repo := repository.NewPreferencesRepository(db)
	urepo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	savedUser, err := urepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	preference := users.Preference{
		UserID:      savedUser.ID,
		EmailEnable: rand.Intn(2) == 0,
		PushEnable:  rand.Intn(2) == 0,
	}
	if _, err := repo.Create(context.Background(), preference); err != nil {
		t.Fatalf("failed to create preference: %v", err)
	}

	cases := []struct {
		desc       string
		preference users.Preference
		err        error
	}{
		{
			desc: "valid preference",
			preference: users.Preference{
				UserID:     savedUser.ID,
				PushEnable: rand.Intn(2) == 0,
			},
			err: nil,
		},
		{
			desc: "invalid preference",
			preference: users.Preference{
				UserID:     invalidID,
				PushEnable: rand.Intn(2) == 0,
			},
			err: errors.New("preference not found"),
		},
		{
			desc: "malformed preference",
			preference: users.Preference{
				UserID:     malformedID,
				PushEnable: rand.Intn(2) == 0,
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			preference, err := repo.UpdatePush(context.Background(), tc.preference)
			switch {
			case err == nil:
				assert.Equal(t, tc.preference.PushEnable, preference.PushEnable)
			default:
				assert.Empty(t, preference)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestDeletePreference(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM preferences")
		require.NoError(t, err)
	})
	repo := repository.NewPreferencesRepository(db)
	urepo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	savedUser, err := urepo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	preference := users.Preference{
		UserID:      savedUser.ID,
		EmailEnable: rand.Intn(2) == 0,
		PushEnable:  rand.Intn(2) == 0,
	}
	if _, err := repo.Create(context.Background(), preference); err != nil {
		t.Fatalf("failed to create preference: %v", err)
	}

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "valid user",
			id:   savedUser.ID,
			err:  nil,
		},
		{
			desc: "invalid user",
			id:   invalidID,
			err:  errors.New("preference not found"),
		},
		{
			desc: "malformed user",
			id:   malformedID,
			err:  errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.Delete(context.Background(), tc.id)
			if err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}
