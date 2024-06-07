// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func ExtractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if strings.HasPrefix(token, "Bearer ") {
		return strings.TrimPrefix(token, "Bearer ")
	}

	return token
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)

		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch {
	case strings.Contains(err.Error(), "not found"):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, errors.New("unauthorized")):
		w.WriteHeader(http.StatusUnauthorized)
	case strings.Contains(err.Error(), "invalid input syntax") ||
		strings.Contains(err.Error(), "insert or update on table") ||
		strings.Contains(err.Error(), "value too long") ||
		strings.Contains(err.Error(), "required") ||
		strings.Contains(err.Error(), "empty"):
		w.WriteHeader(http.StatusBadRequest)
	case strings.Contains(err.Error(), "new row for relation"):
		w.WriteHeader(http.StatusConflict)
	case strings.Contains(err.Error(), "null value") ||
		strings.Contains(err.Error(), "the provided hex string is not a valid ObjectID"):
		w.WriteHeader(http.StatusUnprocessableEntity)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
