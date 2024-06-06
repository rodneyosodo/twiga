// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/rodneyosodo/twiga/internal/postgres"
	"github.com/rodneyosodo/twiga/users"
)

var _ users.FeedRepository = (*feedRepository)(nil)

type feedRepository struct {
	postgres.Database
}

func NewFeedRepository(db postgres.Database) *feedRepository {
	return &feedRepository{db}
}

func (r *feedRepository) Create(ctx context.Context, feed users.Feed) error {
	query := `INSERT INTO feeds (user_id, post_id)
			VALUES (:user_id, :post_id)`

	result, err := r.NamedExecContext(ctx, query, feed)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("could not create feed")
	}

	return nil
}

func (r *feedRepository) RetrieveAll(ctx context.Context, page users.Page) (fpage users.FeedPage, err error) {
	filter := ""
	filters := []string{}
	if page.UserID != "" {
		filters = append(filters, fmt.Sprintf("user_id = '%s'", page.UserID))
	}
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " AND "))
	}

	query := fmt.Sprintf(`SELECT * FROM feeds %s ORDER BY created_at DESC LIMIT :limit OFFSET :offset`, filter)
	dPage := users.FollowingsPage{
		Page: page,
	}

	rows, err := r.NamedQueryContext(ctx, query, dPage)
	if err != nil {
		return users.FeedPage{}, err
	}
	defer rows.Close()

	items := make([]users.Feed, 0)
	for rows.Next() {
		var feed users.Feed
		if err := rows.StructScan(&feed); err != nil {
			return users.FeedPage{}, err
		}

		items = append(items, feed)
	}

	totalQuery := `SELECT COUNT(*) FROM feeds`
	if filter != "" {
		totalQuery = fmt.Sprintf(`SELECT COUNT(*) FROM feeds %s`, filter)
	}

	total, err := postgres.Total(ctx, r.Database, totalQuery, dPage)
	if err != nil {
		return users.FeedPage{}, err
	}

	return users.FeedPage{
		Page: users.Page{
			Limit:  page.Limit,
			Offset: page.Offset,
			Total:  total,
		},
		Feeds: items,
	}, nil
}
