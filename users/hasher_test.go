// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package users_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rodneyosodo/twiga/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPaossword(t *testing.T) {
	cases := []struct {
		desc     string
		password string
		err      error
	}{
		{
			desc:     "empty password",
			password: "",
			err:      errors.New("empty password"),
		},
		{
			desc:     "valid password",
			password: "password",
			err:      nil,
		},
		{
			desc:     "long password",
			password: strings.Repeat("a", 100),
			err:      errors.New("bcrypt: password length exceeds 72 bytes"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			hash, err := users.HashPassword(tc.password)
			switch err {
			case nil:
				assert.NotEmpty(t, hash)
				err := users.ComparePassword(tc.password, hash)
				assert.NoError(t, err)
			default:
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	hash, err := users.HashPassword("password")
	require.NoError(t, err)

	cases := []struct {
		desc     string
		password string
		hash     string
		err      error
	}{
		{
			desc:     "empty password",
			password: "",
			hash:     "hash",
			err:      errors.New("empty password"),
		},
		{
			desc:     "empty hash",
			password: "password",
			hash:     "",
			err:      errors.New("empty hash"),
		},
		{
			desc:     "invalid password",
			password: "12345678",
			hash:     "hash",
			err:      errors.New("compare hash and password failed"),
		},
		{
			desc:     "valid password",
			password: "password",
			hash:     hash,
			err:      nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := users.ComparePassword(tc.password, tc.hash)
			switch err {
			case nil:
				assert.NoError(t, err)
			default:
				assert.ErrorContains(t, err, tc.err.Error())
			}
		})
	}
}
