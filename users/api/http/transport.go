// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	iapi "github.com/rodneyosodo/twiga/internal/api"
	"github.com/rodneyosodo/twiga/users"
	"github.com/rodneyosodo/twiga/users/api"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func MakeHandler(svc users.Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(iapi.EncodeError),
	}

	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
		r.Post("/", otelhttp.NewHandler(kithttp.NewServer(
			api.CreateUserEndpoint(svc),
			decodeCreateUserRequest,
			iapi.EncodeResponse,
			opts...,
		), "create_user").ServeHTTP)
		r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
			api.GetUsersEndpoint(svc),
			decodeGetUsersRequest,
			iapi.EncodeResponse,
			opts...,
		), "get_user").ServeHTTP)

		r.Patch("/password", otelhttp.NewHandler(kithttp.NewServer(
			api.UpdatePasswordEndpoint(svc),
			decodeUpdateUserPasswordRequest,
			iapi.EncodeResponse,
			opts...,
		), "update_user_password").ServeHTTP)

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
				api.GetUserEndpoint(svc),
				decodeEntityRequest,
				iapi.EncodeResponse,
				opts...,
			), "get_user").ServeHTTP)
			r.Put("/", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdateUserEndpoint(svc, api.AllFields),
				decodeUpdateUserRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user").ServeHTTP)
			r.Delete("/", otelhttp.NewHandler(kithttp.NewServer(
				api.DeleteUserEndpoint(svc),
				decodeEntityRequest,
				iapi.EncodeResponse,
				opts...,
			), "delete_user").ServeHTTP)

			r.Patch("/usermame", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdateUserEndpoint(svc, api.UsernameField),
				decodeUpdateUserRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_username").ServeHTTP)
			r.Patch("/email", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdateUserEndpoint(svc, api.EmailField),
				decodeUpdateUserRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_email").ServeHTTP)
			r.Patch("/bio", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdateUserEndpoint(svc, api.BioField),
				decodeUpdateUserRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_bio").ServeHTTP)
			r.Patch("/picture", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdateUserEndpoint(svc, api.PictureField),
				decodeUpdateUserRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_picture").ServeHTTP)
			r.Patch("/preferences", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdateUserEndpoint(svc, api.PreferencesField),
				decodeUpdateUserRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_preferences").ServeHTTP)

			r.Post("/follow", otelhttp.NewHandler(kithttp.NewServer(
				api.FollowEndpoint(svc),
				decodeFollowRequest,
				iapi.EncodeResponse,
				opts...,
			), "follow_user").ServeHTTP)
			r.Post("/unfollow", otelhttp.NewHandler(kithttp.NewServer(
				api.FollowEndpoint(svc),
				decodeFollowRequest,
				iapi.EncodeResponse,
				opts...,
			), "unfollow_user").ServeHTTP)

			r.Get("/followers", otelhttp.NewHandler(kithttp.NewServer(
				api.GetFollowingsEndpoint(svc),
				decodeGetFollowersRequest,
				iapi.EncodeResponse,
				opts...,
			), "get_user_followers").ServeHTTP)
			r.Get("/following", otelhttp.NewHandler(kithttp.NewServer(
				api.GetFollowingsEndpoint(svc),
				decodeGetFollowingsRequest,
				iapi.EncodeResponse,
				opts...,
			), "get_user_following").ServeHTTP)
		})

		r.Route("/token", func(r chi.Router) {
			r.Post("/issue", otelhttp.NewHandler(kithttp.NewServer(
				api.IssueTokenEndpoint(svc),
				decodeIssueTokenRequest,
				iapi.EncodeResponse,
				opts...,
			), "get_token").ServeHTTP)
			r.Post("/refresh", otelhttp.NewHandler(kithttp.NewServer(
				api.RefreshTokenEndpoint(svc),
				decodeRefreshTokenRequest,
				iapi.EncodeResponse,
				opts...,
			), "refresh_token").ServeHTTP)
		})

		r.Route("/preferences", func(r chi.Router) {
			r.Post("/", otelhttp.NewHandler(kithttp.NewServer(
				api.CreatePreferenceEndpoint(svc),
				decodeUpdatePreferenceRequest,
				iapi.EncodeResponse,
				opts...,
			), "create_user_preferences").ServeHTTP)
			r.Get("/", otelhttp.NewHandler(kithttp.NewServer(
				api.GetUserPreferenceEndpoint(svc),
				decodeCreatePreferenceRequest,
				iapi.EncodeResponse,
				opts...,
			), "get_user_preferences").ServeHTTP)
			r.Put("/", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdatePreferenceEndpoint(svc, api.AllFields),
				decodeUpdatePreferenceRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_preferences").ServeHTTP)
			r.Patch("/email", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdatePreferenceEndpoint(svc, api.EmailField),
				decodeUpdatePreferenceRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_email_preferences").ServeHTTP)
			r.Patch("/push", otelhttp.NewHandler(kithttp.NewServer(
				api.UpdatePreferenceEndpoint(svc, api.PushField),
				decodeUpdatePreferenceRequest,
				iapi.EncodeResponse,
				opts...,
			), "update_user_push_preferences").ServeHTTP)
		})

		r.Get("/feed", otelhttp.NewHandler(kithttp.NewServer(
			api.GetFeedEndpoint(svc),
			decodeGetUsersRequest,
			iapi.EncodeResponse,
			opts...,
		), "get_feed").ServeHTTP)
	})

	return r
}

func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.EntitiesReq
	req.Token = iapi.ExtractToken(r)
	offset := r.URL.Query().Get("offset")
	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}
	req.Page.Offset = uint64(intOffset)

	limit := r.URL.Query().Get("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	req.Page.Limit = uint64(intLimit)

	return req, nil
}

func decodeEntityRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.EntityReq
	req.ID = chi.URLParam(r, "userID")
	req.Token = iapi.ExtractToken(r)

	return req, nil
}

func decodeUpdateUserPasswordRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.UpdatePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Token = iapi.ExtractToken(r)

	return req, nil
}

func decodeUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Token = iapi.ExtractToken(r)
	req.User.ID = chi.URLParam(r, "userID")

	return req, nil
}

func decodeIssueTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.IssueTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.RefreshTokenReq
	req.Token = iapi.ExtractToken(r)

	return req, nil
}

func decodeCreatePreferenceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.CreatePreferenceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Token = iapi.ExtractToken(r)

	return req, nil
}

func decodeUpdatePreferenceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.UpdatePreferenceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	req.Token = iapi.ExtractToken(r)

	return req, nil
}

func decodeFollowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.FollowReq
	req.Token = iapi.ExtractToken(r)
	req.ID = chi.URLParam(r, "userID")

	return req, nil
}

func decodeGetFollowersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.EntitiesReq
	req.Token = iapi.ExtractToken(r)
	offset := r.URL.Query().Get("offset")
	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}
	req.Page.Offset = uint64(intOffset)

	limit := r.URL.Query().Get("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	req.Page.Limit = uint64(intLimit)

	req.Page.FollowerID = chi.URLParam(r, "userID")

	return req, nil
}

func decodeGetFollowingsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.EntitiesReq
	req.Token = iapi.ExtractToken(r)
	offset := r.URL.Query().Get("offset")
	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}
	req.Page.Offset = uint64(intOffset)

	limit := r.URL.Query().Get("limit")
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	req.Page.Limit = uint64(intLimit)

	req.Page.FolloweeID = chi.URLParam(r, "userID")

	return req, nil
}
