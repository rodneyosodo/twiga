// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const cost int = 10

var (
	errHashPassword    = errors.New("generate hash from password failed")
	errComparePassword = errors.New("compare hash and password failed")
)

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return "", errors.Join(errHashPassword, err)
	}

	return string(hash), nil
}

func ComparePassword(plain, hashed string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		return errors.Join(errComparePassword, err)
	}

	return nil
}
