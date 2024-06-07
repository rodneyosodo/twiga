// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package users

import (
	"context"
	"encoding/json"
	"time"
)

type Page struct {
	Total      uint64 `db:"total"       json:"total"`
	Offset     uint64 `db:"offset"      json:"offset"`
	Limit      uint64 `db:"limit"       json:"limit"`
	FollowerID string `db:"follower_id" json:"follower_id"`
	FolloweeID string `db:"followee_id" json:"followee_id"`
	UserID     string `db:"user_id"     json:"user_id"`
}

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	PictureURL  string    `json:"picture_url"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Preferences []string  `json:"preferences"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u User) KV() (string, string, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return "", "", err
	}

	return u.ID, string(data), nil
}

func (u User) FromKV(key, value string) error {
	return json.Unmarshal([]byte(value), &u)
}

type UsersPage struct {
	Page
	Users []User `json:"users"`
}

type Preference struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	EmailEnable bool      `json:"email_enabled"`
	PushEnable  bool      `json:"push_enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p Preference) KV() (string, string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", "", err
	}

	return p.ID, string(data), nil
}

func (p Preference) FromKV(key, value string) error {
	return json.Unmarshal([]byte(value), &p)
}

type PreferencesPage struct {
	Page
	Preferences []Preference `json:"preferences"`
}

type Following struct {
	ID         string    `db:"id"          json:"id"`
	FollowerID string    `db:"follower_id" json:"follower_id"`
	FolloweeID string    `db:"followee_id" json:"followee_id"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}

type FollowingsPage struct {
	Page
	Followings []Following `json:"followings"`
}

type Feed struct {
	ID        string    `db:"id"         json:"id"`
	UserID    string    `db:"user_id"    json:"user_id"`
	PostID    string    `db:"post_id"    json:"post_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type FeedPage struct {
	Page
	Feeds []Feed `json:"feeds"`
}

func (page UsersPage) MarshalJSON() ([]byte, error) {
	type Alias UsersPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Users == nil {
		a.Users = make([]User, 0)
	}

	return json.Marshal(a)
}

func (page PreferencesPage) MarshalJSON() ([]byte, error) {
	type Alias PreferencesPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Preferences == nil {
		a.Preferences = make([]Preference, 0)
	}

	return json.Marshal(a)
}

func (page FeedPage) MarshalJSON() ([]byte, error) {
	type Alias FeedPage
	a := struct {
		Alias
	}{
		Alias: Alias(page),
	}

	if a.Feeds == nil {
		a.Feeds = make([]Feed, 0)
	}

	return json.Marshal(a)
}

//go:generate mockery --name UsersRepository --output=./mocks --filename users.go --quiet
type UsersRepository interface { //nolint:interfacebloat
	Create(ctx context.Context, user User) (User, error)
	RetrieveByID(ctx context.Context, id string) (User, error)
	RetrieveAll(ctx context.Context, page Page) (UsersPage, error)
	Update(ctx context.Context, user User) (User, error)
	UpdateUsername(ctx context.Context, user User) (User, error)
	UpdatePassword(ctx context.Context, user User) (User, error)
	UpdateEmail(ctx context.Context, user User) (User, error)
	UpdateBio(ctx context.Context, user User) (User, error)
	UpdatePictureURL(ctx context.Context, user User) (User, error)
	UpdatePreferences(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
}

//go:generate mockery --name PreferencesRepository --output=./mocks --filename preference.go --quiet
type PreferencesRepository interface {
	Create(ctx context.Context, preference Preference) (Preference, error)
	RetrieveByUserID(ctx context.Context, userID string) (Preference, error)
	RetrieveAll(ctx context.Context, page Page) (PreferencesPage, error)
	Update(ctx context.Context, preference Preference) (Preference, error)
	UpdateEmail(ctx context.Context, preference Preference) (Preference, error)
	UpdatePush(ctx context.Context, preference Preference) (Preference, error)
	Delete(ctx context.Context, userID string) error
}

//go:generate mockery --name FollowingRepository --output=./mocks --filename following.go --quiet
type FollowingRepository interface {
	Create(ctx context.Context, following Following) (Following, error)
	RetrieveAll(ctx context.Context, page Page) (FollowingsPage, error)
	Delete(ctx context.Context, following Following) error
}

//go:generate mockery --name FeedRepository --output=./mocks --filename feed.go --quiet
type FeedRepository interface {
	Create(ctx context.Context, feed Feed) error
	RetrieveAll(ctx context.Context, page Page) (FeedPage, error)
}

//go:generate mockery --name Service --output=./mocks --filename service.go --quiet
type Service interface { //nolint:interfacebloat
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUsers(ctx context.Context, page Page) (UsersPage, error)
	UpdateUser(ctx context.Context, user User) (User, error)
	UpdateUserUsername(ctx context.Context, user User) (User, error)
	UpdateUserPassword(ctx context.Context, id, oldPassword, currentPassowrd string) (User, error)
	UpdateUserEmail(ctx context.Context, user User) (User, error)
	UpdateUserBio(ctx context.Context, user User) (User, error)
	UpdateUserPictureURL(ctx context.Context, user User) (User, error)
	UpdateUserPreferences(ctx context.Context, user User) (User, error)
	DeleteUser(ctx context.Context, id string) error

	CreatePreferences(ctx context.Context, preference Preference) (Preference, error)
	GetPreferencesByUserID(ctx context.Context, userID string) (Preference, error)
	GetPreferences(ctx context.Context, page Page) (PreferencesPage, error)
	UpdatePreferences(ctx context.Context, preference Preference) (Preference, error)
	UpdateEmailPreferences(ctx context.Context, preference Preference) (Preference, error)
	UpdatePushPreferences(ctx context.Context, preference Preference) (Preference, error)
	DeletePreferences(ctx context.Context, userID string) error

	CreateFollower(ctx context.Context, following Following) (Following, error)
	GetUserFollowings(ctx context.Context, page Page) (FollowingsPage, error)
	DeleteFollower(ctx context.Context, following Following) error

	CreateFeed(ctx context.Context, feed Feed) error
	GetUserFeed(ctx context.Context, page Page) (FeedPage, error)
}
