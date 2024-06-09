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
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/rodneyosodo/twiga/users/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const svcName = "proto.UsersService"

var _ proto.UsersServiceClient = (*client)(nil)

type client struct {
	timeout            time.Duration
	getUserByID        endpoint.Endpoint
	getUserPreferences endpoint.Endpoint
	getUserFollowers   endpoint.Endpoint
	createFeed         endpoint.Endpoint
	identifyUser       endpoint.Endpoint
}

func NewClient(conn *grpc.ClientConn, timeout time.Duration) proto.UsersServiceClient {
	return &client{
		getUserByID: kitgrpc.NewClient(
			conn,
			svcName,
			"GetUserByID",
			encodeGetUserByIDRequest,
			decodeGetUserByIDResponse,
			proto.GetUserByIDResponse{},
		).Endpoint(),

		getUserPreferences: kitgrpc.NewClient(
			conn,
			svcName,
			"GetUserPreferences",
			encodeGetUserPreferencesRequest,
			decodeGetUserPreferencesResponse,
			proto.GetUserPreferencesResponse{},
		).Endpoint(),

		getUserFollowers: kitgrpc.NewClient(
			conn,
			svcName,
			"GetUserFollowers",
			encodeGetUserFollowersRequest,
			decodeGetUserFollowersResponse,
			proto.GetUserFollowersResponse{},
		).Endpoint(),

		createFeed: kitgrpc.NewClient(
			conn,
			svcName,
			"CreateFeed",
			encodeCreateFeedRequest,
			decodeCreateFeedResponse,
			proto.CreateFeedResponse{},
		).Endpoint(),

		identifyUser: kitgrpc.NewClient(
			conn,
			svcName,
			"IdentifyUser",
			encodeIdentifyUserRequest,
			decodeIdentifyUserResponse,
			proto.IdentifyUserResponse{},
		).Endpoint(),

		timeout: timeout,
	}
}

func (c *client) GetUserByID(ctx context.Context, in *proto.GetUserByIDRequest, _ ...grpc.CallOption) (*proto.GetUserByIDResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := c.getUserByID(ctx, in)
	if err != nil {
		return nil, decodeError(err)
	}

	response, ok := res.(*proto.GetUserByIDResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return response, nil
}

func encodeGetUserByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.GetUserByIDRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}

	return req, nil
}

func decodeGetUserByIDResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(*proto.GetUserByIDResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return res, nil
}

func (c *client) GetUserPreferences(ctx context.Context, in *proto.GetUserPreferencesRequest, _ ...grpc.CallOption) (*proto.GetUserPreferencesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := c.getUserPreferences(ctx, in)
	if err != nil {
		return nil, decodeError(err)
	}

	response, ok := res.(*proto.GetUserPreferencesResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return response, nil
}

func encodeGetUserPreferencesRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.GetUserPreferencesRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}

	return req, nil
}

func decodeGetUserPreferencesResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(*proto.GetUserPreferencesResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return res, nil
}

func (c *client) GetUserFollowers(ctx context.Context, in *proto.GetUserFollowersRequest, _ ...grpc.CallOption) (*proto.GetUserFollowersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := c.getUserFollowers(ctx, in)
	if err != nil {
		return nil, decodeError(err)
	}

	response, ok := res.(*proto.GetUserFollowersResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return response, nil
}

func encodeGetUserFollowersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.GetUserFollowersRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}

	return req, nil
}

func decodeGetUserFollowersResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(*proto.GetUserFollowersResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return res, nil
}

func (c *client) CreateFeed(ctx context.Context, in *proto.CreateFeedRequest, _ ...grpc.CallOption) (*proto.CreateFeedResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := c.createFeed(ctx, in)
	if err != nil {
		return nil, decodeError(err)
	}

	response, ok := res.(*proto.CreateFeedResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return response, nil
}

func encodeCreateFeedRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.CreateFeedRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}

	return req, nil
}

func decodeCreateFeedResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(*proto.CreateFeedResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return res, nil
}

func (c *client) IdentifyUser(ctx context.Context, in *proto.IdentifyUserRequest, _ ...grpc.CallOption) (*proto.IdentifyUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := c.identifyUser(ctx, in)
	if err != nil {
		return nil, decodeError(err)
	}

	response, ok := res.(*proto.IdentifyUserResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return response, nil
}

func encodeIdentifyUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proto.IdentifyUserRequest)
	if !ok {
		return nil, errors.New("invalid request")
	}

	return req, nil
}

func decodeIdentifyUserResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res, ok := grpcRes.(*proto.IdentifyUserResponse)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return res, nil
}

func decodeError(err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Unauthenticated, codes.PermissionDenied, codes.InvalidArgument, codes.FailedPrecondition, codes.NotFound, codes.AlreadyExists, codes.OK:
			return errors.New(st.Message())
		default:
			return errors.Join(fmt.Errorf("unexpected gRPC status: %s (status code:%v)", st.Code().String(), st.Code()), errors.New(st.Message()))
		}
	}

	return err
}
