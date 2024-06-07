// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package notifications

import (
	"context"
	"errors"

	"github.com/rodneyosodo/twiga/users/proto"
)

var _ Service = (*service)(nil)

type service struct {
	repo  Repository
	users proto.UsersServiceClient
}

func NewService(repo Repository, users proto.UsersServiceClient) Service {
	return &service{
		repo:  repo,
		users: users,
	}
}

func (s *service) CreateNotification(ctx context.Context, token string, notification Notification) (Notification, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Notification{}, err
	}
	notification.UserID = userID

	return s.repo.CreateNotification(ctx, notification)
}

func (s *service) RetrieveNotification(ctx context.Context, token string, id string) (Notification, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Notification{}, err
	}
	n, err := s.repo.RetrieveNotification(ctx, id)
	if err != nil {
		return Notification{}, err
	}
	if n.UserID != userID {
		return Notification{}, errors.New("unauthorized")
	}

	return n, nil
}

func (s *service) RetrieveAllNotifications(ctx context.Context, token string, page Page) (NotificationsPage, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return NotificationsPage{}, err
	}
	page.UserID = userID

	return s.repo.RetrieveAllNotifications(ctx, page)
}

func (s *service) ReadNotification(ctx context.Context, token string, id string) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}

	return s.repo.ReadNotification(ctx, userID, id)
}

func (s *service) ReadAllNotifications(ctx context.Context, token string, page Page) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	page.UserID = userID

	return s.repo.ReadAllNotifications(ctx, page)
}

func (s *service) DeleteNotification(ctx context.Context, token string, id string) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	n, err := s.repo.RetrieveNotification(ctx, id)
	if err != nil {
		return err
	}
	if n.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteNotification(ctx, id)
}

func (s *service) CreateSetting(ctx context.Context, token string, setting Setting) (Setting, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Setting{}, err
	}
	setting.UserID = userID

	return s.repo.CreateSetting(ctx, setting)
}

func (s *service) RetrieveSetting(ctx context.Context, token string, id string) (Setting, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return Setting{}, err
	}
	saved, err := s.repo.RetrieveSetting(ctx, id)
	if err != nil {
		return Setting{}, err
	}
	if saved.UserID != userID {
		return Setting{}, errors.New("unauthorized")
	}

	return saved, nil
}

func (s *service) RetrieveAllSettings(ctx context.Context, token string, page Page) (SettingsPage, error) {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return SettingsPage{}, err
	}
	page.UserID = userID

	return s.repo.RetrieveAllSettings(ctx, page)
}

func (s *service) UpdateSetting(ctx context.Context, token string, setting Setting) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	saved, err := s.repo.RetrieveSetting(ctx, setting.ID)
	if err != nil {
		return err
	}
	if saved.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.UpdateSetting(ctx, setting)
}

func (s *service) UpdateEmailSetting(ctx context.Context, token string, id string, isEnabled bool) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	saved, err := s.repo.RetrieveSetting(ctx, id)
	if err != nil {
		return err
	}
	if saved.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.UpdateEmailSetting(ctx, id, isEnabled)
}

func (s *service) UpdatePushSetting(ctx context.Context, token string, id string, isEnabled bool) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	saved, err := s.repo.RetrieveSetting(ctx, id)
	if err != nil {
		return err
	}
	if saved.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.UpdatePushSetting(ctx, id, isEnabled)
}

func (s *service) DeleteSetting(ctx context.Context, token string, id string) error {
	userID, err := s.IdentifyUser(ctx, token)
	if err != nil {
		return err
	}
	saved, err := s.repo.RetrieveSetting(ctx, id)
	if err != nil {
		return err
	}
	if saved.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteSetting(ctx, id)
}

func (s *service) IdentifyUser(ctx context.Context, token string) (string, error) {
	resp, err := s.users.IdentifyUser(ctx, &proto.IdentifyUserRequest{Token: token})
	if err != nil {
		return "", err
	}

	return resp.GetId(), nil
}
