// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package jwt_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	tjwt "github.com/rodneyosodo/twiga/users/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const secret = "longsecretkey"

func newToken(t *testing.T, expired bool, issuer bool) string {
	id, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("generating new UUID expected to succeed: %s", err))
	userID, err := uuid.NewV4()
	require.Nil(t, err, fmt.Sprintf("generating new UUID expected to succeed: %s", err))
	expiration := time.Now().Add(1 * time.Minute)
	if expired {
		expiration = time.Now().Add(-1 * time.Minute)
	}
	issuerValue := "twiga"
	if issuer {
		issuerValue = "invalid"
	}
	token, err := jwt.NewBuilder().
		Issuer(issuerValue).
		Subject(userID.String()).
		JwtID(id.String()).
		IssuedAt(time.Now()).
		Expiration(expiration).
		Build()
	require.Nil(t, err, fmt.Sprintf("building token expected to succeed: %s", err))
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, []byte(secret)))
	require.Nil(t, err, fmt.Sprintf("signing token expected to succeed: %s", err))

	return string(signed)
}

func TestIssue(t *testing.T) {
	tokenizer := tjwt.NewTokenizer(secret, time.Minute)

	cases := []struct {
		desc   string
		userID string
		err    error
	}{
		{
			desc:   "issue new token",
			userID: uuid.Must(uuid.NewV4()).String(),
			err:    nil,
		},
		{
			desc: "issue with empty user",
			err:  errors.New("user id is required"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			tkn, err := tokenizer.Issue(tc.userID)
			switch tc.err {
			case nil:
				assert.NotEmpty(t, tkn)
				userID, err := tokenizer.Validate(tkn)
				assert.Nil(t, err)
				assert.Equal(t, tc.userID, userID)
			default:
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tokenizer := tjwt.NewTokenizer(secret, time.Minute)

	cases := []struct {
		desc string
		tkn  string
		err  error
	}{
		{
			desc: "valid token",
			tkn:  newToken(t, false, false),
			err:  nil,
		},
		{
			desc: "expired token",
			tkn:  newToken(t, true, false),
			err:  errors.New("\"exp\" not satisfied"),
		},
		{
			desc: "invalid issuer",
			tkn:  newToken(t, false, true),
			err:  errors.New("invalid token issuer value"),
		},
		{
			desc: "invalid token",
			tkn:  "invalid",
			err:  errors.New("invalid JWT"),
		},
		{
			desc: "empty token",
			tkn:  "",
			err:  errors.New("token is required"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			userID, err := tokenizer.Validate(tc.tkn)
			switch tc.err {
			case nil:
				assert.NotEmpty(t, userID)
			default:
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}
