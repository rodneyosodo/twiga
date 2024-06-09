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
	SVC   SVC
}

func (req EntityReq) validate() error {
	if req.ID == "" {
		return errors.New("id is required")
	}

	if req.Token == "" && req.SVC == HTTP {
		return errors.New("token is required")
	}

	return nil
}

type EntitiesReq struct {
	Token string
	Page  users.Page
	SVC   SVC
}

func (req EntitiesReq) validate() error {
	if req.Token == "" && req.SVC == HTTP {
		return errors.New("token is required")
	}

	return nil
}

type UpdateUserReq struct {
	Token string
	users.User
}

func (req UpdateUserReq) validate(field entityField) error {
	if req.ID == "" {
		return errors.New("id is required")
	}

	if req.Token == "" {
		return errors.New("token is required")
	}

	if req.Username == "" && field == UsernameField {
		return errors.New("username is required")
	}

	if req.Email == "" && field == EmailField {
		return errors.New("email is required")
	}

	if req.Bio == "" && field == BioField {
		return errors.New("bio is required")
	}
	if req.PictureURL == "" && field == PictureField {
		return errors.New("picture_url is required")
	}
	if req.Preferences == nil && field == PreferencesField {
		return errors.New("preferences is required")
	}

	return nil
}

type UpdatePasswordReq struct {
	Token           string
	OldPassword     string `json:"old_password"`
	CurrentPassowrd string `json:"current_password"`
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
	users.User
}

func (req IssueTokenReq) validate() error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Password == "" {
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
	Token string
	users.Preference
}

func (req CreatePreferenceReq) validate() error {
	if req.Token == "" {
		return errors.New("token is required")
	}

	return nil
}

type SVC uint8

const (
	HTTP SVC = iota
	GRPC
)

type GetPreferenceReq struct {
	Token string
	ID    string
	SVC   SVC
}

func (req GetPreferenceReq) validate() error {
	if req.Token == "" && req.SVC == HTTP {
		return errors.New("token is required")
	}

	return nil
}

type UpdatePreferenceReq struct {
	Token string
	users.Preference
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
	users.Feed
}

func (req CreateFeedReq) validate() error {
	return nil
}
