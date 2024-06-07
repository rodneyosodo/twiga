// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package api

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/rodneyosodo/twiga/users"
)

type entityField string

const (
	UsernameField    entityField = "username"
	EmailField       entityField = "email"
	BioField         entityField = "bio"
	PictureField     entityField = "picture"
	PreferencesField entityField = "preferences"
	PushField        entityField = "push"
	AllFields        entityField = "all"
)

func CreateUserEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateUserReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.CreateUser(ctx, req.User)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func GetUserEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EntityReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.GetUserByID(ctx, req.ID, req.Token)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func GetUsersEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EntitiesReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		users, err := svc.GetUsers(ctx, req.Token, req.Page)
		if err != nil {
			return nil, err
		}

		return users, nil
	}
}

func UpdateUserEndpoint(svc users.Service, field entityField) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateUserReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(field); err != nil {
			return nil, err
		}

		var user users.User
		var err error
		switch field {
		case UsernameField:
			user, err = svc.UpdateUserUsername(ctx, req.Token, req.User)
		case EmailField:
			user, err = svc.UpdateUserEmail(ctx, req.Token, req.User)
		case BioField:
			user, err = svc.UpdateUserBio(ctx, req.Token, req.User)
		case PictureField:
			user, err = svc.UpdateUserPictureURL(ctx, req.Token, req.User)
		case PreferencesField:
			user, err = svc.UpdateUserPreferences(ctx, req.Token, req.User)
		default:
			user, err = svc.UpdateUser(ctx, req.Token, req.User)
		}
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func UpdatePasswordEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdatePasswordReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		user, err := svc.UpdateUserPassword(ctx, req.Token, req.OldPassword, req.CurrentPassowrd)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

func DeleteUserEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EntityReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.DeleteUser(ctx, req.Token, req.ID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func IssueTokenEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(IssueTokenReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		token, err := svc.IssueToken(ctx, req.User)
		if err != nil {
			return nil, err
		}

		return token, nil
	}
}

func RefreshTokenEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(RefreshTokenReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		token, err := svc.RefreshToken(ctx, req.Token)
		if err != nil {
			return nil, err
		}

		return token, nil
	}
}

func CreatePreferenceEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreatePreferenceReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		preference, err := svc.CreatePreferences(ctx, req.Token, req.Preference)
		if err != nil {
			return nil, err
		}

		return preference, nil
	}
}

func GetUserPreferenceEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EntityReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		preference, err := svc.GetPreferencesByUserID(ctx, req.Token)
		if err != nil {
			return nil, err
		}

		return preference, nil
	}
}

func UpdatePreferenceEndpoint(svc users.Service, field entityField) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdatePreferenceReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		var preference users.Preference
		var err error

		switch field {
		case EmailField:
			preference, err = svc.UpdateEmailPreferences(ctx, req.Token, req.Preference)
		case PushField:
			preference, err = svc.UpdatePushPreferences(ctx, req.Token, req.Preference)
		default:
			preference, err = svc.UpdatePreferences(ctx, req.Token, req.Preference)
		}
		if err != nil {
			return nil, err
		}

		return preference, nil
	}
}

func DeletePreferenceEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EntityReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		err := svc.DeletePreferences(ctx, req.Token)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func FollowEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FollowReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		following := users.Following{
			FolloweeID: req.ID,
		}
		if _, err := svc.CreateFollower(ctx, req.Token, following); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func UnfollowEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(FollowReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		following := users.Following{
			FolloweeID: req.ID,
		}
		if err := svc.DeleteFollower(ctx, req.Token, following); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func GetFollowingsEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EntitiesReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		followings, err := svc.GetUserFollowings(ctx, req.Token, req.Page)
		if err != nil {
			return nil, err
		}

		return followings, nil
	}
}

func CreateFeedEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateFeedReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		if err := svc.CreateFeed(ctx, req.Feed); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func GetFeedEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EntitiesReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		feed, err := svc.GetUserFeed(ctx, req.Token, req.Page)
		if err != nil {
			return nil, err
		}

		return feed, nil
	}
}

func IdentifyUserEndpoint(svc users.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(RefreshTokenReq)
		if !ok {
			return nil, errors.New("invalid request")
		}

		if err := req.validate(); err != nil {
			return nil, err
		}

		userID, err := svc.IdentifyUser(ctx, req.Token)
		if err != nil {
			return nil, err
		}

		return userID, nil
	}
}
