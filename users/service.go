// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package users

import (
	"context"
	"strings"

	"github.com/0x6flab/namegenerator"
	"github.com/gofrs/uuid"
	"github.com/rodneyosodo/twiga/internal/cache"
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
	cacher          cache.Cache
}

func NewService(usersRepo UsersRepository, preferencesRepo PreferencesRepository, followingRepo FollowingRepository, feedRepo FeedRepository) Service {
	return &service{
		usersRepo:       usersRepo,
		preferencesRepo: preferencesRepo,
		followingRepo:   followingRepo,
		feedRepo:        feedRepo,
	}
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

	return s.usersRepo.Create(ctx, user)
}

func (s *service) GetUserByID(ctx context.Context, id string) (User, error) {
	return s.usersRepo.RetrieveByID(ctx, id)
}

func (s *service) GetUsers(ctx context.Context, page Page) (UsersPage, error) {
	return s.usersRepo.RetrieveAll(ctx, page)
}

func (s *service) UpdateUser(ctx context.Context, user User) (User, error) {
	return s.usersRepo.Update(ctx, user)
}

func (s *service) UpdateUserUsername(ctx context.Context, user User) (User, error) {
	return s.usersRepo.UpdateUsername(ctx, user)
}

func (s *service) UpdateUserPassword(ctx context.Context, id, oldPassword, currentPassowrd string) (User, error) {
	user, err := s.usersRepo.RetrieveByID(ctx, id)
	if err != nil {
		return User{}, err
	}
	if err := ComparePassword(oldPassword, user.Password); err != nil {
		return User{}, err
	}

	user.Password, err = HashPassword(currentPassowrd)
	if err != nil {
		return User{}, err
	}

	return s.usersRepo.UpdatePassword(ctx, user)
}

func (s *service) UpdateUserEmail(ctx context.Context, user User) (User, error) {
	return s.usersRepo.UpdateEmail(ctx, user)
}

func (s *service) UpdateUserBio(ctx context.Context, user User) (User, error) {
	return s.usersRepo.UpdateBio(ctx, user)
}

func (s *service) UpdateUserPictureURL(ctx context.Context, user User) (User, error) {
	return s.usersRepo.UpdatePictureURL(ctx, user)
}

func (s *service) UpdateUserPreferences(ctx context.Context, user User) (User, error) {
	return s.usersRepo.UpdatePreferences(ctx, user)
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	return s.usersRepo.Delete(ctx, id)
}

func (s *service) CreatePreferences(ctx context.Context, preference Preference) (Preference, error) {
	return s.preferencesRepo.Create(ctx, preference)
}

func (s *service) GetPreferencesByUserID(ctx context.Context, userID string) (Preference, error) {
	return s.preferencesRepo.RetrieveByUserID(ctx, userID)
}

func (s *service) GetPreferences(ctx context.Context, page Page) (PreferencesPage, error) {
	return s.preferencesRepo.RetrieveAll(ctx, page)
}

func (s *service) UpdatePreferences(ctx context.Context, preference Preference) (Preference, error) {
	return s.preferencesRepo.Update(ctx, preference)
}

func (s *service) UpdateEmailPreferences(ctx context.Context, preference Preference) (Preference, error) {
	return s.preferencesRepo.UpdateEmail(ctx, preference)
}

func (s *service) UpdatePushPreferences(ctx context.Context, preference Preference) (Preference, error) {
	return s.preferencesRepo.UpdatePush(ctx, preference)
}

func (s *service) DeletePreferences(ctx context.Context, userID string) error {
	return s.preferencesRepo.Delete(ctx, userID)
}

func (s *service) CreateFollower(ctx context.Context, following Following) (Following, error) {
	return s.followingRepo.Create(ctx, following)
}

func (s *service) GetUserFollowings(ctx context.Context, page Page) (FollowingsPage, error) {
	return s.followingRepo.RetrieveAll(ctx, page)
}

func (s *service) DeleteFollower(ctx context.Context, following Following) error {
	return s.followingRepo.Delete(ctx, following)
}

func (s *service) CreateFeed(ctx context.Context, feed Feed) error {
	return s.feedRepo.Create(ctx, feed)
}

func (s *service) GetUserFeed(ctx context.Context, page Page) (FeedPage, error) {
	return s.feedRepo.RetrieveAll(ctx, page)
}
