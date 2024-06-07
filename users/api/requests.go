// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package api

import (
	"errors"

	"github.com/rodneyosodo/twiga/users"
)

type CreateUserReq struct {
	users.User
}

func (req CreateUserReq) validate() error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

type EntityReq struct {
	ID    string
	Token string
}

func (req EntityReq) validate() error {
	if req.ID == "" {
		return errors.New("id is required")
	}

	if req.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type EntitiesReq struct {
	Token string
	Page  users.Page
}

func (req EntitiesReq) validate() error {
	if req.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type UpdateUserReq struct {
	Token string
	User  users.User
}

func (req UpdateUserReq) validate(field entityField) error {
	if req.User.ID == "" {
		return errors.New("id is required")
	}

	if req.Token == "" {
		return errors.New("token is required")
	}

	if req.User.Username == "" && field == UsernameField {
		return errors.New("username is required")
	}

	if req.User.Email == "" && field == EmailField {
		return errors.New("email is required")
	}

	if req.User.Bio == "" && field == BioField {
		return errors.New("bio is required")
	}
	if req.User.PictureURL == "" && field == PictureField {
		return errors.New("picture_url is required")
	}
	if req.User.Preferences == nil && field == PreferencesField {
		return errors.New("preferences is required")
	}

	return nil
}

type UpdatePasswordReq struct {
	Token           string
	OldPassword     string
	CurrentPassowrd string
}

func (req UpdatePasswordReq) validate() error {
	if req.Token == "" {
		return errors.New("token is required")
	}

	if req.OldPassword == "" {
		return errors.New("old password is required")
	}

	if req.CurrentPassowrd == "" {
		return errors.New("current password is required")
	}

	return nil
}

type IssueTokenReq struct {
	User users.User
}

func (req IssueTokenReq) validate() error {
	if req.User.Email == "" {
		return errors.New("email is required")
	}

	if req.User.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

type RefreshTokenReq struct {
	Token string
}

func (req RefreshTokenReq) validate() error {
	if req.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type CreatePreferenceReq struct {
	Token      string
	Preference users.Preference
}

func (req CreatePreferenceReq) validate() error {
	if req.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type UpdatePreferenceReq struct {
	Token      string
	Preference users.Preference
}

func (req UpdatePreferenceReq) validate() error {
	if req.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type FollowReq struct {
	Token string
	ID    string
}

func (req FollowReq) validate() error {
	if req.Token == "" {
		return errors.New("token is required")
	}

	if req.ID == "" {
		return errors.New("id is required")
	}

	return nil
}

type CreateFeedReq struct {
	Feed users.Feed
}

func (req CreateFeedReq) validate() error {
	return nil
}
