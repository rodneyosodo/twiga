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
	"time"

	"github.com/jackc/pgtype"
	"github.com/rodneyosodo/twiga/internal/postgres"
	"github.com/rodneyosodo/twiga/users"
)

var _ users.PreferencesRepository = (*pRepository)(nil)

type pRepository struct {
	postgres.Database
}

func NewPreferencesRepository(db postgres.Database) *pRepository {
	return &pRepository{db}
}

func (r *pRepository) Create(ctx context.Context, preference users.Preference) (users.Preference, error) {
	query := `INSERT INTO preferences (user_id, email_enabled, push_enabled)
			VALUES (:user_id, :email_enabled, :push_enabled)
			RETURNING *`
	dPreference, err := toDBPreference(preference)
	if err != nil {
		return users.Preference{}, err
	}

	rows, err := r.NamedQueryContext(ctx, query, dPreference)
	if err != nil {
		return users.Preference{}, err
	}
	defer rows.Close()

	dPreference = dbPreference{}
	if rows.Next() {
		if err := rows.StructScan(&dPreference); err != nil {
			return users.Preference{}, err
		}
	}

	return fromDBPreference(dPreference), nil
}

func (r *pRepository) RetrieveByUserID(ctx context.Context, userID string) (users.Preference, error) {
	query := `SELECT * FROM preferences WHERE user_id = :user_id`
	dPreference := dbPreference{
		UserID: userID,
	}

	rows, err := r.NamedQueryContext(ctx, query, dPreference)
	if err != nil {
		return users.Preference{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var dPreference dbPreference
		if err := rows.StructScan(&dPreference); err != nil {
			return users.Preference{}, err
		}

		return fromDBPreference(dPreference), nil
	}

	return users.Preference{}, errors.New("preference not found")
}

func (r *pRepository) RetrieveAll(ctx context.Context, page users.Page) (ppage users.PreferencesPage, err error) {
	query := `SELECT * FROM preferences ORDER BY created_at DESC LIMIT :limit OFFSET :offset`
	dPage := users.PreferencesPage{
		Page: page,
	}

	rows, err := r.NamedQueryContext(ctx, query, dPage)
	if err != nil {
		return users.PreferencesPage{}, err
	}
	defer rows.Close()

	items := make([]users.Preference, 0)
	for rows.Next() {
		var dPreference dbPreference
		if err := rows.StructScan(&dPreference); err != nil {
			return users.PreferencesPage{}, err
		}

		items = append(items, fromDBPreference(dPreference))
	}

	totalQuery := `SELECT COUNT(*) FROM preferences`

	total, err := postgres.Total(ctx, r.Database, totalQuery, dPage)
	if err != nil {
		return users.PreferencesPage{}, err
	}

	return users.PreferencesPage{
		Page: users.Page{
			Limit:  page.Limit,
			Offset: page.Offset,
			Total:  total,
		},
		Preferences: items,
	}, nil
}

func (r *pRepository) Update(ctx context.Context, preference users.Preference) (users.Preference, error) {
	query := `UPDATE preferences SET email_enabled = :email_enabled, push_enabled = :push_enabled, updated_at = CURRENT_TIMESTAMP
			WHERE user_id = :user_id
			RETURNING *`

	return r.update(ctx, query, preference)
}

func (r *pRepository) UpdateEmail(ctx context.Context, preference users.Preference) (users.Preference, error) {
	query := `UPDATE preferences SET email_enabled = :email_enabled, updated_at = CURRENT_TIMESTAMP
			WHERE user_id = :user_id
			RETURNING *`

	return r.update(ctx, query, preference)
}

func (r *pRepository) UpdatePush(ctx context.Context, preference users.Preference) (users.Preference, error) {
	query := `UPDATE preferences SET push_enabled = :push_enabled, updated_at = CURRENT_TIMESTAMP
			WHERE user_id = :user_id
			RETURNING *`

	return r.update(ctx, query, preference)
}

func (r *pRepository) update(ctx context.Context, query string, preference users.Preference) (users.Preference, error) {
	dPreference, err := toDBPreference(preference)
	if err != nil {
		return users.Preference{}, err
	}

	rows, err := r.NamedQueryContext(ctx, query, dPreference)
	if err != nil {
		return users.Preference{}, err
	}
	defer rows.Close()

	dPreference = dbPreference{}
	if rows.Next() {
		if err := rows.StructScan(&dPreference); err != nil {
			return users.Preference{}, err
		}

		return fromDBPreference(dPreference), nil
	}

	return users.Preference{}, errors.New("preference not found")
}

func (r *pRepository) Delete(ctx context.Context, userID string) error {
	q := "DELETE FROM preferences WHERE user_id = $1;"

	result, err := r.ExecContext(ctx, q, userID)
	if err != nil {
		return errors.Join(errors.New("could not delete user preference"), err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("preference not found")
	}

	return nil
}

type dbPreference struct {
	ID          string      `db:"id"`
	UserID      string      `db:"user_id"`
	EmailEnable pgtype.Bool `db:"email_enabled"`
	PushEnable  pgtype.Bool `db:"push_enabled"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

func toDBPreference(preference users.Preference) (dbPreference, error) {
	emailEnabled := pgtype.Bool{}
	if err := emailEnabled.Set(preference.EmailEnable); err != nil {
		return dbPreference{}, err
	}

	pushEnabled := pgtype.Bool{}
	if err := pushEnabled.Set(preference.PushEnable); err != nil {
		return dbPreference{}, err
	}

	return dbPreference{
		ID:          preference.ID,
		UserID:      preference.UserID,
		EmailEnable: emailEnabled,
		PushEnable:  pushEnabled,
		CreatedAt:   preference.CreatedAt,
		UpdatedAt:   preference.UpdatedAt,
	}, nil
}

func fromDBPreference(preference dbPreference) users.Preference {
	return users.Preference{
		ID:          preference.ID,
		UserID:      preference.UserID,
		EmailEnable: preference.EmailEnable.Bool,
		PushEnable:  preference.PushEnable.Bool,
		CreatedAt:   preference.CreatedAt,
		UpdatedAt:   preference.UpdatedAt,
	}
}
