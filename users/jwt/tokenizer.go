// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rodneyosodo/twiga/users"
)

const issuer = "twiga"

type tokenizer struct {
	secret  []byte
	expires time.Duration
}

func NewTokenizer(secret string, expires time.Duration) users.Tokenizer {
	return &tokenizer{
		secret:  []byte(secret),
		expires: expires,
	}
}

func (t *tokenizer) Issue(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user id is required")
	}
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	token, err := jwt.NewBuilder().
		Issuer(issuer).
		Subject(userID).
		JwtID(id.String()).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(t.expires)).
		Build()
	if err != nil {
		return "", err
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS512, t.secret))
	if err != nil {
		return "", err
	}

	return string(signedToken), nil
}

func (t *tokenizer) Validate(token string) (string, error) {
	if token == "" {
		return "", errors.New("token is required")
	}

	tkn, err := jwt.Parse(
		[]byte(token),
		jwt.WithValidate(true),
		jwt.WithKey(jwa.HS512, t.secret),
	)
	if err != nil {
		return "", err
	}
	validator := jwt.ValidatorFunc(func(_ context.Context, jwtToken jwt.Token) jwt.ValidationError {
		if jwtToken.Issuer() != issuer {
			return jwt.NewValidationError(errors.New("invalid token issuer value"))
		}

		return nil
	})
	if err := jwt.Validate(tkn, jwt.WithValidator(validator)); err != nil {
		return "", err
	}

	return tkn.Subject(), nil
}
