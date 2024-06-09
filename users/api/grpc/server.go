// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package grpc

import (
	"context"
	"errors"
	"strings"
	"time"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/rodneyosodo/twiga/users"
	"github.com/rodneyosodo/twiga/users/api"
	"github.com/rodneyosodo/twiga/users/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ proto.UsersServiceServer = (*server)(nil)

type server struct {
	proto.UnimplementedUsersServiceServer
	getUserByID        kitgrpc.Handler
	getUserPreferences kitgrpc.Handler
	getUserFollowers   kitgrpc.Handler
	createFeed         kitgrpc.Handler
	identifyUser       kitgrpc.Handler
}

func NewServer(svc users.Service) proto.UsersServiceServer {
	return &server{
		getUserByID: kitgrpc.NewServer(
			api.GetUserEndpoint(svc),
			decodeGetUserByIDRequest,
			encodeGetUserByIDResponse,
		),
		getUserPreferences: kitgrpc.NewServer(
			api.GetUserPreferenceEndpoint(svc),
			decodeGetUserPreferencesRequest,
			encodeGetUserPreferencesResponse,
		),
		getUserFollowers: kitgrpc.NewServer(
			api.GetFollowingsEndpoint(svc),
			decodeGetUserFollowersRequest,
			encodeGetUserFollowersResponse,
		),
		createFeed: kitgrpc.NewServer(
			api.CreateFeedEndpoint(svc),
			decodeCreateFeedRequest,
			encodeCreateFeedResponse,
		),
		identifyUser: kitgrpc.NewServer(
			api.IdentifyUserEndpoint(svc),
			decodeIdentifyUserRequest,
			encodeIdentifyUserResponse,
		),
	}
}

func (s *server) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserByIDResponse, error) {
	_, res, err := s.getUserByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	resp, ok := res.(*proto.GetUserByIDResponse)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return resp, nil
}

func decodeGetUserByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.GetUserByIDRequest)
	if !ok {
		return nil, encodeError(errors.New("invalid request"))
	}

	return api.EntityReq{
		ID:  req.GetId(),
		SVC: api.GRPC,
	}, nil
}

func encodeGetUserByIDResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(users.User)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return &proto.GetUserByIDResponse{
		Id:             res.ID,
		Username:       res.Username,
		DisplayName:    res.DisplayName,
		Bio:            res.Bio,
		ProfilePicture: res.PictureURL,
		CreatedAt:      res.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      res.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *server) GetUserPreferences(ctx context.Context, req *proto.GetUserPreferencesRequest) (*proto.GetUserPreferencesResponse, error) {
	_, res, err := s.getUserPreferences.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	resp, ok := res.(*proto.GetUserPreferencesResponse)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return resp, nil
}

func decodeGetUserPreferencesRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.GetUserPreferencesRequest)
	if !ok {
		return nil, encodeError(errors.New("invalid request"))
	}

	return api.GetPreferenceReq{
		SVC: api.GRPC,
		ID:  req.GetId(),
	}, nil
}

func encodeGetUserPreferencesResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(users.Preference)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return &proto.GetUserPreferencesResponse{
		EmailNotifications: res.EmailEnable,
		PushNotifications:  res.PushEnable,
	}, nil
}

func (s *server) GetUserFollowers(ctx context.Context, req *proto.GetUserFollowersRequest) (*proto.GetUserFollowersResponse, error) {
	_, res, err := s.getUserFollowers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	resp, ok := res.(*proto.GetUserFollowersResponse)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return resp, nil
}

func decodeGetUserFollowersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.GetUserFollowersRequest)
	if !ok {
		return nil, encodeError(errors.New("invalid request"))
	}

	return api.EntitiesReq{
		SVC: api.GRPC,
		Page: users.Page{
			Offset: req.GetOffset(),
			Limit:  req.GetLimit(),
			UserID: req.GetId(),
		},
	}, nil
}

func encodeGetUserFollowersResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(users.FollowingsPage)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	followings := make([]*proto.Following, 0, len(res.Followings))
	for _, following := range res.Followings {
		followings = append(followings, &proto.Following{
			Id:         following.ID,
			FollowerId: following.FollowerID,
			FolloweeId: following.FolloweeID,
		})
	}

	return &proto.GetUserFollowersResponse{
		Followings: followings,
		Total:      res.Total,
		Offset:     res.Offset,
		Limit:      res.Limit,
	}, nil
}

func (s *server) CreateFeed(ctx context.Context, req *proto.CreateFeedRequest) (*proto.CreateFeedResponse, error) {
	_, res, err := s.createFeed.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	resp, ok := res.(*proto.CreateFeedResponse)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return resp, nil
}

func decodeCreateFeedRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.CreateFeedRequest)
	if !ok {
		return nil, encodeError(errors.New("invalid request"))
	}

	return api.CreateFeedReq{
		Feed: users.Feed{
			UserID: req.GetUserId(),
			PostID: req.GetPostId(),
		},
	}, nil
}

func encodeCreateFeedResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(map[string]string)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return &proto.CreateFeedResponse{
		Message: res["message"],
	}, nil
}

func (s *server) IdentifyUser(ctx context.Context, req *proto.IdentifyUserRequest) (*proto.IdentifyUserResponse, error) {
	_, res, err := s.identifyUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}
	resp, ok := res.(*proto.IdentifyUserResponse)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return resp, nil
}

func decodeIdentifyUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.IdentifyUserRequest)
	if !ok {
		return nil, encodeError(errors.New("invalid request"))
	}

	return api.RefreshTokenReq{
		Token: req.GetToken(),
	}, nil
}

func encodeIdentifyUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(string)
	if !ok {
		return nil, encodeError(errors.New("invalid response"))
	}

	return &proto.IdentifyUserResponse{Id: res}, nil
}

func encodeError(err error) error {
	switch {
	case strings.Contains(err.Error(), "not found"):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, errors.New("unauthorized")):
		return status.Error(codes.PermissionDenied, err.Error())
	case strings.Contains(err.Error(), "invalid input syntax") ||
		strings.Contains(err.Error(), "insert or update on table") ||
		strings.Contains(err.Error(), "value too long") ||
		strings.Contains(err.Error(), "required") ||
		strings.Contains(err.Error(), "empty"):
		return status.Error(codes.InvalidArgument, err.Error())
	case strings.Contains(err.Error(), "new row for relation"):
		return status.Error(codes.AlreadyExists, err.Error())
	case strings.Contains(err.Error(), "null value") ||
		strings.Contains(err.Error(), "the provided hex string is not a valid ObjectID"):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
