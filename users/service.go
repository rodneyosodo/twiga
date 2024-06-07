// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package users

import (
	"context"
	"errors"
	"strings"

	"github.com/0x6flab/namegenerator"
	"github.com/gofrs/uuid"
)

var (
	namegen   = namegenerator.NewGenerator()
	defAvatar = "https://ui-avatars.com/api/?name="
)

type service struct {
	usersRepo       UsersRepository
	preferencesRepo PreferencesRepository
	followingRepo   FollowingRepository
	feedRepo        FeedRepository
	tokenizer       Tokenizer
}

func NewService(usersRepo UsersRepository, preferencesRepo PreferencesRepository, followingRepo FollowingRepository, feedRepo FeedRepository, tokenizer Tokenizer) Service {
	return &service{
		usersRepo:       usersRepo,
		preferencesRepo: preferencesRepo,
		followingRepo:   followingRepo,
		feedRepo:        feedRepo,
		tokenizer:       tokenizer,
	}
}

func (s *service) IssueToken(ctx context.Context, user User) (string, error) {
	saved, err := s.usersRepo.RetrieveByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := ComparePassword(user.Password, saved.Password); err != nil {
		return "", err
	}

	return s.tokenizer.Issue(saved.ID)
}

func (s *service) RefreshToken(ctx context.Context, token string) (string, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return "", err
	}
	if _, err := s.usersRepo.RetrieveByID(ctx, userID); err != nil {
		return "", err
	}

	return s.tokenizer.Issue(userID)
}

func (s *service) CreateUser(ctx context.Context, user User) (User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}
	user.ID = id.String()

	name := namegen.Generate()
	if user.Username == "" {
		user.Username = strings.ToLower(name)
	}
	if user.DisplayName == "" {
		user.DisplayName = name
	}
	if user.PictureURL == "" {
		user.PictureURL = defAvatar + user.DisplayName
	}
	if user.Email == "" {
		return User{}, errors.New("email is required")
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return User{}, err
	}
	user.Password = hashedPassword

	return s.usersRepo.Create(ctx, user)
}

func (s *service) GetUserByID(ctx context.Context, token string, id string) (User, error) {
	return s.usersRepo.RetrieveByID(ctx, id)
}

func (s *service) GetUsers(ctx context.Context, token string, page Page) (UsersPage, error) {
	if _, err := s.tokenizer.Validate(token); err != nil {
		return UsersPage{}, err
	}

	return s.usersRepo.RetrieveAll(ctx, page)
}

func (s *service) UpdateUser(ctx context.Context, token string, user User) (User, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return User{}, err
	}
	saved, err := s.usersRepo.RetrieveByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	if saved.ID != user.ID {
		return User{}, errors.New("unauthorized")
	}

	return s.usersRepo.Update(ctx, user)
}

func (s *service) UpdateUserUsername(ctx context.Context, token string, user User) (User, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return User{}, err
	}
	saved, err := s.usersRepo.RetrieveByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	if saved.ID != user.ID {
		return User{}, errors.New("unauthorized")
	}
	return s.usersRepo.UpdateUsername(ctx, user)
}

func (s *service) UpdateUserPassword(ctx context.Context, token string, oldPassword, currentPassowrd string) error {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return err
	}

	user, err := s.usersRepo.RetrieveByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.ID != userID {
		return errors.New("unauthorized")
	}

	user, err = s.usersRepo.RetrieveByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if err := ComparePassword(oldPassword, user.Password); err != nil {
		return err
	}

	user.Password, err = HashPassword(currentPassowrd)
	if err != nil {
		return err
	}

	return s.usersRepo.UpdatePassword(ctx, user)
}

func (s *service) UpdateUserEmail(ctx context.Context, token string, user User) (User, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return User{}, err
	}
	saved, err := s.usersRepo.RetrieveByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	if saved.ID != user.ID {
		return User{}, errors.New("unauthorized")
	}

	return s.usersRepo.UpdateEmail(ctx, user)
}

func (s *service) UpdateUserBio(ctx context.Context, token string, user User) (User, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return User{}, err
	}
	saved, err := s.usersRepo.RetrieveByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	if saved.ID != user.ID {
		return User{}, errors.New("unauthorized")
	}

	return s.usersRepo.UpdateBio(ctx, user)
}

func (s *service) UpdateUserPictureURL(ctx context.Context, token string, user User) (User, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return User{}, err
	}
	saved, err := s.usersRepo.RetrieveByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	if saved.ID != user.ID {
		return User{}, errors.New("unauthorized")
	}

	return s.usersRepo.UpdatePictureURL(ctx, user)
}

func (s *service) UpdateUserPreferences(ctx context.Context, token string, user User) (User, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return User{}, err
	}
	saved, err := s.usersRepo.RetrieveByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	if saved.ID != user.ID {
		return User{}, errors.New("unauthorized")
	}

	return s.usersRepo.UpdatePreferences(ctx, user)
}

func (s *service) DeleteUser(ctx context.Context, token string, id string) error {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return err
	}
	if userID != id {
		return errors.New("unauthorized")
	}

	return s.usersRepo.Delete(ctx, id)
}

func (s *service) CreatePreferences(ctx context.Context, token string, preference Preference) (Preference, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return Preference{}, err
	}
	preference.UserID = userID

	return s.preferencesRepo.Create(ctx, preference)
}

func (s *service) GetPreferencesByUserID(ctx context.Context, token, id string) (Preference, error) {
	if token != "" {
		userID, err := s.tokenizer.Validate(token)
		if err != nil {
			return Preference{}, err
		}

		id = userID
	}

	return s.preferencesRepo.RetrieveByUserID(ctx, id)
}

func (s *service) GetPreferences(ctx context.Context, token string, page Page) (PreferencesPage, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return PreferencesPage{}, err
	}
	page.UserID = userID

	return s.preferencesRepo.RetrieveAll(ctx, page)
}

func (s *service) UpdatePreferences(ctx context.Context, token string, preference Preference) (Preference, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return Preference{}, err
	}
	preference.UserID = userID

	return s.preferencesRepo.Update(ctx, preference)
}

func (s *service) UpdateEmailPreferences(ctx context.Context, token string, preference Preference) (Preference, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return Preference{}, err
	}
	preference.UserID = userID

	return s.preferencesRepo.UpdateEmail(ctx, preference)
}

func (s *service) UpdatePushPreferences(ctx context.Context, token string, preference Preference) (Preference, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return Preference{}, err
	}
	preference.UserID = userID

	return s.preferencesRepo.UpdatePush(ctx, preference)
}

func (s *service) DeletePreferences(ctx context.Context, token string) error {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return err
	}

	return s.preferencesRepo.Delete(ctx, userID)
}

func (s *service) CreateFollower(ctx context.Context, token string, following Following) (Following, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return Following{}, err
	}
	following.FollowerID = userID

	return s.followingRepo.Create(ctx, following)
}

func (s *service) GetUserFollowings(ctx context.Context, token string, page Page) (FollowingsPage, error) {
	if _, err := s.tokenizer.Validate(token); err != nil {
		return FollowingsPage{}, err
	}

	return s.followingRepo.RetrieveAll(ctx, page)
}

func (s *service) DeleteFollower(ctx context.Context, token string, following Following) error {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return err
	}
	following.FollowerID = userID

	return s.followingRepo.Delete(ctx, following)
}

func (s *service) CreateFeed(ctx context.Context, feed Feed) error {
	return s.feedRepo.Create(ctx, feed)
}

func (s *service) GetUserFeed(ctx context.Context, token string, page Page) (FeedPage, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return FeedPage{}, err
	}
	page.UserID = userID

	return s.feedRepo.RetrieveAll(ctx, page)
}

func (s *service) IdentifyUser(ctx context.Context, token string) (string, error) {
	userID, err := s.tokenizer.Validate(token)
	if err != nil {
		return "", err
	}

	return userID, nil
}
