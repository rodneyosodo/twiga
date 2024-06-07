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
	"github.com/rodneyosodo/twiga/users"
	"github.com/rodneyosodo/twiga/users/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	invalidID   = uuid.Must(uuid.NewV4()).String()
	malformedID = strings.Repeat("a", 100)
	namegen     = namegenerator.NewGenerator()
)

func TestCreateUser(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				Username:    namegen.Generate(),
				DisplayName: namegen.Generate(),
				Bio:         strings.Repeat("a", 100),
				PictureURL:  "https://example.com" + namegen.Generate(),
				Email:       namegen.Generate() + "@example.com",
				Password:    "password",
				Preferences: []string{"test"},
			},
			err: nil,
		},
		{
			desc: "empty user",
			user: users.User{},
			err:  errors.New("null value"),
		},
		{
			desc: "malformed user",
			user: users.User{
				Username:    strings.Repeat("a", 1025),
				DisplayName: strings.Repeat("a", 1025),
				Bio:         strings.Repeat("a", 100),
				PictureURL:  "https://example.com" + namegen.Generate(),
				Email:       namegen.Generate() + "@example.com",
				Password:    "password",
				Preferences: []string{"test"},
			},
			err: errors.New("value too long"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.Create(context.Background(), tc.user)
			switch {
			case err == nil:
				assert.Equal(t, tc.user.Username, user.Username)
				assert.Equal(t, tc.user.DisplayName, user.DisplayName)
				assert.Equal(t, tc.user.Bio, user.Bio)
				assert.Equal(t, tc.user.PictureURL, user.PictureURL)
				assert.Equal(t, tc.user.Email, user.Email)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestRetrieveUserByID(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		id   string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			id:   saved.ID,
			user: saved,
			err:  nil,
		},
		{
			desc: "invalid user",
			id:   invalidID,
			user: users.User{},
			err:  errors.New("user not found"),
		},
		{
			desc: "malformed user",
			id:   malformedID,
			user: users.User{},
			err:  errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.RetrieveByID(context.Background(), tc.id)
			if err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			}
			assert.Equal(t, tc.user, user)
		})
	}
}

func TestRetrieveUserByEmail(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc  string
		email string
		user  users.User
		err   error
	}{
		{
			desc:  "valid user",
			email: saved.Email,
			user:  saved,
			err:   nil,
		},
		{
			desc:  "invalid user",
			email: namegen.Generate() + "@example.com",
			user:  users.User{},
			err:   errors.New("user not found"),
		},
		{
			desc:  "malformed user",
			email: malformedID,
			user:  users.User{},
			err:   errors.New("user not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.RetrieveByEmail(context.Background(), tc.email)
			if err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			}
			tc.user.Password = user.Password
			assert.Equal(t, tc.user, user)
		})
	}
}

func TestRetrieveAllUsers(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	num := uint64(10)
	savedUsers := make([]users.User, num)

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
		saved, err := repo.Create(context.Background(), user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
		savedUsers[i] = saved
	}

	cases := []struct { //nolint:dupl
		desc      string
		page      users.Page
		usersPage users.UsersPage
		err       error
	}{
		{
			desc: "valid page",
			page: users.Page{
				Offset: 0,
				Limit:  10,
			},
			usersPage: users.UsersPage{
				Page: users.Page{
					Offset: 0,
					Limit:  10,
					Total:  uint64(len(savedUsers)),
				},
				Users: savedUsers,
			},
			err: nil,
		},
		{
			desc: "big offset",
			page: users.Page{
				Offset: num + 1,
				Limit:  10,
			},
			usersPage: users.UsersPage{
				Page: users.Page{
					Offset: num + 1,
					Limit:  10,
					Total:  uint64(len(savedUsers)),
				},
				Users: []users.User{},
			},
			err: nil,
		},
		{
			desc: "big limit",
			page: users.Page{
				Offset: 0,
				Limit:  num + 1,
			},
			usersPage: users.UsersPage{
				Page: users.Page{
					Offset: 0,
					Limit:  num + 1,
					Total:  uint64(len(savedUsers)),
				},
				Users: savedUsers,
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			usersPage, err := repo.RetrieveAll(context.Background(), tc.page)
			switch {
			case err == nil:
				assert.Equal(t, tc.usersPage.Offset, usersPage.Offset)
				assert.Equal(t, tc.usersPage.Limit, usersPage.Limit)
				assert.Equal(t, tc.usersPage.Total, usersPage.Total)
				assert.ElementsMatch(t, tc.usersPage.Users, usersPage.Users)
			default:
				assert.Empty(t, usersPage)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				ID:          saved.ID,
				Username:    namegen.Generate(),
				DisplayName: namegen.Generate(),
				Bio:         strings.Repeat("a", 100),
				PictureURL:  "https://example.com" + namegen.Generate(),
				Email:       namegen.Generate() + "@example.com",
				Preferences: []string{"test"},
			},
			err: nil,
		},
		{
			desc: "invalid user",
			user: users.User{
				ID:          invalidID,
				Username:    namegen.Generate(),
				DisplayName: namegen.Generate(),
				Bio:         strings.Repeat("a", 100),
			},
			err: errors.New("user not found"),
		},
		{
			desc: "malformed user",
			user: users.User{
				ID:          malformedID,
				Username:    namegen.Generate(),
				DisplayName: namegen.Generate(),
				Bio:         strings.Repeat("a", 100),
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.Update(context.Background(), tc.user)
			switch {
			case err == nil:
				assert.Equal(t, tc.user.Username, user.Username)
				assert.Equal(t, tc.user.DisplayName, user.DisplayName)
				assert.Equal(t, tc.user.Bio, user.Bio)
				assert.Equal(t, tc.user.PictureURL, user.PictureURL)
				assert.Equal(t, tc.user.Email, user.Email)
				assert.Equal(t, tc.user.Preferences, user.Preferences)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdateUsername(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				ID:       saved.ID,
				Username: namegen.Generate(),
			},
			err: nil,
		},
		{
			desc: "invalid user",
			user: users.User{
				ID:       invalidID,
				Username: namegen.Generate(),
			},
			err: errors.New("user not found"),
		},
		{
			desc: "malformed user",
			user: users.User{
				ID:       malformedID,
				Username: namegen.Generate(),
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.UpdateUsername(context.Background(), tc.user)
			switch {
			case err == nil:
				assert.Equal(t, tc.user.Username, user.Username)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				ID:       saved.ID,
				Password: "newpassword",
			},
			err: nil,
		},
		{
			desc: "invalid user",
			user: users.User{
				ID:       invalidID,
				Password: "newpassword",
			},
			err: errors.New("user not found"),
		},
		{
			desc: "malformed user",
			user: users.User{
				ID:       malformedID,
				Password: "newpassword",
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := repo.UpdatePassword(context.Background(), tc.user)
			switch {
			case err == nil:
				user, err := repo.RetrieveByEmail(context.Background(), saved.Email)
				require.NoError(t, err)
				assert.Equal(t, tc.user.Password, user.Password)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdateEmail(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				ID:    saved.ID,
				Email: namegen.Generate() + "@example.com",
			},
			err: nil,
		},
		{
			desc: "invalid user",
			user: users.User{
				ID:    invalidID,
				Email: namegen.Generate() + "@example.com",
			},
			err: errors.New("user not found"),
		},
		{
			desc: "malformed user",
			user: users.User{
				ID:    malformedID,
				Email: namegen.Generate() + "@example.com",
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.UpdateEmail(context.Background(), tc.user)
			switch {
			case err == nil:
				assert.Equal(t, tc.user.Email, user.Email)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdateBio(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				ID:  saved.ID,
				Bio: strings.Repeat("b", 100),
			},
			err: nil,
		},
		{
			desc: "invalid user",
			user: users.User{
				ID:  invalidID,
				Bio: strings.Repeat("b", 100),
			},
			err: errors.New("user not found"),
		},
		{
			desc: "malformed user",
			user: users.User{
				ID:  malformedID,
				Bio: strings.Repeat("b", 100),
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.UpdateBio(context.Background(), tc.user)
			switch {
			case err == nil:
				assert.Equal(t, tc.user.Bio, user.Bio)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdatePictureURL(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				ID:         saved.ID,
				PictureURL: "https://example.com" + namegen.Generate(),
			},
			err: nil,
		},
		{
			desc: "invalid user",
			user: users.User{
				ID:         invalidID,
				PictureURL: "https://example.com" + namegen.Generate(),
			},
			err: errors.New("user not found"),
		},
		{
			desc: "malformed user",
			user: users.User{
				ID:         malformedID,
				PictureURL: "https://example.com" + namegen.Generate(),
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.UpdatePictureURL(context.Background(), tc.user)
			switch {
			case err == nil:
				assert.Equal(t, tc.user.PictureURL, user.PictureURL)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestUpdatePreferences(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    namegen.Generate(),
		DisplayName: namegen.Generate(),
		Bio:         strings.Repeat("a", 100),
		PictureURL:  "https://example.com" + namegen.Generate(),
		Email:       namegen.Generate() + "@example.com",
		Password:    "password",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		user users.User
		err  error
	}{
		{
			desc: "valid user",
			user: users.User{
				ID:          saved.ID,
				Preferences: []string{namegen.Generate()},
			},
			err: nil,
		},
		{
			desc: "invalid user",
			user: users.User{
				ID:          invalidID,
				Preferences: []string{namegen.Generate()},
			},
			err: errors.New("user not found"),
		},
		{
			desc: "malformed user",
			user: users.User{
				ID:          malformedID,
				Preferences: []string{namegen.Generate()},
			},
			err: errors.New("invalid input syntax"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			user, err := repo.UpdatePreferences(context.Background(), tc.user)
			switch {
			case err == nil:
				assert.Equal(t, tc.user.Preferences, user.Preferences)
			default:
				assert.Empty(t, user)
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	t.Cleanup(func() {
		_, err := db.Exec("DELETE FROM users")
		require.NoError(t, err)
	})
	repo := repository.NewUsersRepository(db)

	user := users.User{
		Username:    "test",
		DisplayName: "test",
		Bio:         "test",
		PictureURL:  "test",
		Email:       "test",
		Password:    "test",
		Preferences: []string{"test"},
	}
	saved, err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "valid user",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "invalid user",
			id:   invalidID,
			err:  errors.New("user not found"),
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

func generateUser() users.User {
	name := namegen.Generate()

	return users.User{
		Username:    strings.ToLower(name),
		DisplayName: name,
		Bio:         strings.Repeat(name, 100),
		PictureURL:  "https://example.com" + name,
		Email:       name + "@example.com",
		Password:    strings.ToLower(name),
		Preferences: namegen.GenerateMultiple(5),
	}
}
