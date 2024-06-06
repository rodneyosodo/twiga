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

var _ users.FollowingRepository = (*fRepository)(nil)

type fRepository struct {
	postgres.Database
}

func NewFollowingRepository(db postgres.Database) *fRepository {
	return &fRepository{db}
}

func (r *fRepository) Create(ctx context.Context, follower users.Following) (users.Following, error) {
	query := `INSERT INTO followers (follower_id, followee_id)
			VALUES (:follower_id, :followee_id)
			RETURNING *`
	rows, err := r.NamedQueryContext(ctx, query, follower)
	if err != nil {
		return users.Following{}, err
	}
	defer rows.Close()

	follower = users.Following{}
	if rows.Next() {
		if err := rows.StructScan(&follower); err != nil {
			return users.Following{}, err
		}
	}

	return follower, nil
}

func (r *fRepository) RetrieveAll(ctx context.Context, page users.Page) (fpage users.FollowingsPage, err error) {
	filter := ""
	filters := []string{}
	if page.FolloweeID != "" {
		filters = append(filters, fmt.Sprintf("follower_id = '%s'", page.FolloweeID))
	}
	if page.FollowerID != "" {
		filters = append(filters, fmt.Sprintf("followee_id = '%s'", page.FollowerID))
	}
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " AND "))
	}

	query := fmt.Sprintf(`SELECT * FROM followers %s ORDER BY created_at DESC LIMIT :limit OFFSET :offset`, filter)
	dPage := users.FollowingsPage{
		Page: page,
	}

	rows, err := r.NamedQueryContext(ctx, query, dPage)
	if err != nil {
		return users.FollowingsPage{}, err
	}
	defer rows.Close()

	items := make([]users.Following, 0)
	for rows.Next() {
		var following users.Following
		if err := rows.StructScan(&following); err != nil {
			return users.FollowingsPage{}, err
		}

		items = append(items, following)
	}

	totalQuery := `SELECT COUNT(*) FROM followers`
	if filter != "" {
		totalQuery = fmt.Sprintf(`SELECT COUNT(*) FROM followers %s`, filter)
	}

	total, err := postgres.Total(ctx, r.Database, totalQuery, dPage)
	if err != nil {
		return users.FollowingsPage{}, err
	}

	return users.FollowingsPage{
		Page: users.Page{
			Limit:  page.Limit,
			Offset: page.Offset,
			Total:  total,
		},
		Followings: items,
	}, nil
}

func (r *fRepository) Delete(ctx context.Context, following users.Following) error {
	q := "DELETE FROM followers WHERE follower_id = $1 AND followee_id = $2"

	result, err := r.ExecContext(ctx, q, following.FollowerID, following.FolloweeID)
	if err != nil {
		return errors.Join(errors.New("could not delete following"), err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("following not found")
	}

	return nil
}
