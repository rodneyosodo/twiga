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

	"github.com/rodneyosodo/twiga/internal/cache"
	"github.com/rodneyosodo/twiga/users/proto"
)

const defLimit = 100

var _ Service = (*service)(nil)

type service struct {
	repo   Repository
	users  proto.UsersServiceClient
	cacher cache.Cacher
}

func NewService(repo Repository, users proto.UsersServiceClient, cacher cache.Cacher) Service {
	return &service{
		repo:   repo,
		users:  users,
		cacher: cacher,
	}
}

func (s *service) CreateNotification(ctx context.Context, notification Notification) (Notification, error) {
	notification, err := s.repo.CreateNotification(ctx, notification)
	if err != nil {
		return Notification{}, err
	}

	if err := s.cacher.Add(ctx, notification.ID, notification); err != nil {
		return Notification{}, errors.New("failed to cache notification")
	}

	return notification, nil
}

func (s *service) IdentifyUser(ctx context.Context, token string) (string, error) {
	resp, err := s.users.IdentifyUser(ctx, &proto.IdentifyUserRequest{Token: token})
	if err != nil {
		return "", err
	}

	req := &proto.GetUserFollowersRequest{Id: resp.GetId(), Offset: 0, Limit: defLimit}
	followers, err := s.users.GetUserFollowers(ctx, req)
	if err != nil {
		return "", err
	}
	for followers.Total > uint64(len(followers.Followings)) {
		req.Offset += defLimit
		resp, err := s.users.GetUserFollowers(ctx, req)
		if err != nil {
			return "", err
		}
		followers.Followings = append(followers.Followings, resp.Followings...)
	}

	userIDs := make([]string, 0, len(followers.Followings))
	for _, f := range followers.Followings {
		userIDs = append(userIDs, f.FolloweeId)
	}

	if err := s.cacher.Add(ctx, resp.GetId(), userIDs); err != nil {
		return "", errors.New("failed to cache followers")
	}

	return resp.GetId(), nil
}

func (s *service) GetNewNotification(ctx context.Context, userID string) Notification {
	ids, ok := s.cacher.Get(ctx, userID).([]interface{})
	if !ok {
		return Notification{}
	}

	for _, id := range ids {
		id, ok := id.(string)
		if !ok {
			continue
		}

		n, ok := s.cacher.Get(ctx, id).(Notification)
		if ok {
			defer s.cacher.Remove(ctx, id)

			return n
		}
	}

	return Notification{}
}

func (s *service) RetrieveNotification(ctx context.Context, token string, id string) (Notification, error) {
	userID, err := s.identifyUser(ctx, token)
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
	userID, err := s.identifyUser(ctx, token)
	if err != nil {
		return NotificationsPage{}, err
	}
	followers, err := s.users.GetUserFollowers(ctx, &proto.GetUserFollowersRequest{Id: userID, Offset: 0, Limit: defLimit})
	if err != nil {
		return NotificationsPage{}, err
	}
	for followers.Total > defLimit {
		resp, err := s.users.GetUserFollowers(ctx, &proto.GetUserFollowersRequest{Id: userID, Offset: defLimit, Limit: defLimit})
		if err != nil {
			return NotificationsPage{}, err
		}
		followers.Followings = append(followers.Followings, resp.Followings...)
		followers.Total += resp.Total
	}

	userIDs := make([]string, 0, len(followers.Followings))
	for _, f := range followers.Followings {
		userIDs = append(userIDs, f.FolloweeId)
	}
	page.IDs = userIDs

	return s.repo.RetrieveAllNotifications(ctx, page)
}

func (s *service) ReadNotification(ctx context.Context, token string, id string) error {
	userID, err := s.identifyUser(ctx, token)
	if err != nil {
		return err
	}

	return s.repo.ReadNotification(ctx, userID, id)
}

func (s *service) ReadAllNotifications(ctx context.Context, token string, page Page) error {
	userID, err := s.identifyUser(ctx, token)
	if err != nil {
		return err
	}
	page.UserID = userID

	return s.repo.ReadAllNotifications(ctx, page)
}

func (s *service) DeleteNotification(ctx context.Context, token string, id string) error {
	userID, err := s.identifyUser(ctx, token)
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

func (s *service) identifyUser(ctx context.Context, token string) (string, error) {
	resp, err := s.users.IdentifyUser(ctx, &proto.IdentifyUserRequest{Token: token})
	if err != nil {
		return "", err
	}

	return resp.GetId(), nil
}
