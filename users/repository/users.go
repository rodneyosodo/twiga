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
	"time"

	"github.com/jackc/pgtype"
	"github.com/rodneyosodo/twiga/internal/postgres"
	"github.com/rodneyosodo/twiga/users"
)

const defReturnCols = "id, username, display_name, bio, picture_url, email, preferences, created_at, updated_at"

var _ users.UsersRepository = (*uRepository)(nil)

type uRepository struct {
	postgres.Database
}

func NewUsersRepository(db postgres.Database) *uRepository {
	return &uRepository{db}
}

func (r *uRepository) Create(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`INSERT INTO users (username, display_name, bio, picture_url, email, password, preferences)
			VALUES (:username, :display_name, :bio, :picture_url, :email, :password, :preferences)
			RETURNING %s`, defReturnCols)
	dUser, err := toDBUser(user)
	if err != nil {
		return users.User{}, err
	}

	rows, err := r.NamedQueryContext(ctx, query, dUser)
	if err != nil {
		return users.User{}, err
	}
	defer rows.Close()

	dUser = dbUser{}
	if rows.Next() {
		if err := rows.StructScan(&dUser); err != nil {
			return users.User{}, err
		}
	}

	return fromDBUser(dUser), nil
}

func (r *uRepository) RetrieveByID(ctx context.Context, id string) (users.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = :id`, defReturnCols)
	dUser := dbUser{
		ID: id,
	}

	rows, err := r.NamedQueryContext(ctx, query, dUser)
	if err != nil {
		return users.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var dUser dbUser
		if err := rows.StructScan(&dUser); err != nil {
			return users.User{}, err
		}

		return fromDBUser(dUser), nil
	}

	return users.User{}, errors.New("user not found")
}

func (r *uRepository) RetrieveByEmail(ctx context.Context, email string) (users.User, error) {
	query := `SELECT * FROM users WHERE email = :email`
	dUser := dbUser{
		Email: email,
	}

	rows, err := r.NamedQueryContext(ctx, query, dUser)
	if err != nil {
		return users.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var dUser dbUser
		if err := rows.StructScan(&dUser); err != nil {
			return users.User{}, err
		}

		return fromDBUser(dUser), nil
	}

	return users.User{}, errors.New("user not found")
}

func (r *uRepository) RetrieveAll(ctx context.Context, page users.Page) (users.UsersPage, error) {
	query := fmt.Sprintf(`SELECT %s FROM users ORDER BY created_at DESC LIMIT :limit OFFSET :offset`, defReturnCols)
	dPage := users.UsersPage{
		Page: page,
	}

	rows, err := r.NamedQueryContext(ctx, query, dPage)
	if err != nil {
		return users.UsersPage{}, err
	}
	defer rows.Close()

	items := make([]users.User, 0)
	for rows.Next() {
		var dUser dbUser
		if err := rows.StructScan(&dUser); err != nil {
			return users.UsersPage{}, err
		}

		items = append(items, fromDBUser(dUser))
	}

	totalQuery := `SELECT COUNT(*) FROM users`

	total, err := postgres.Total(ctx, r.Database, totalQuery, dPage)
	if err != nil {
		return users.UsersPage{}, err
	}

	return users.UsersPage{
		Page: users.Page{
			Limit:  page.Limit,
			Offset: page.Offset,
			Total:  total,
		},
		Users: items,
	}, nil
}

func (r *uRepository) Update(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`UPDATE users SET username = :username, display_name = :display_name, bio = :bio, picture_url = :picture_url, email = :email, preferences = :preferences, updated_at = CURRENT_TIMESTAMP
			WHERE id = :id
			RETURNING %s`, defReturnCols)

	return r.update(ctx, query, user)
}

func (r *uRepository) UpdateUsername(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`UPDATE users SET username = :username, updated_at = CURRENT_TIMESTAMP
			WHERE id = :id
			RETURNING %s`, defReturnCols)

	return r.update(ctx, query, user)
}

func (r *uRepository) UpdatePassword(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`UPDATE users SET password = :password, updated_at = CURRENT_TIMESTAMP
			WHERE id = :id
			RETURNING %s`, defReturnCols)

	return r.update(ctx, query, user)
}

func (r *uRepository) UpdateEmail(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`UPDATE users SET email = :email, updated_at = CURRENT_TIMESTAMP
			WHERE id = :id
			RETURNING %s`, defReturnCols)

	return r.update(ctx, query, user)
}

func (r *uRepository) UpdateBio(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`UPDATE users SET bio = :bio, updated_at = CURRENT_TIMESTAMP
			WHERE id = :id
			RETURNING %s`, defReturnCols)

	return r.update(ctx, query, user)
}

func (r *uRepository) UpdatePictureURL(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`UPDATE users SET picture_url = :picture_url, updated_at = CURRENT_TIMESTAMP
			WHERE id = :id
			RETURNING %s`, defReturnCols)

	return r.update(ctx, query, user)
}

func (r *uRepository) UpdatePreferences(ctx context.Context, user users.User) (users.User, error) {
	query := fmt.Sprintf(`UPDATE users SET preferences = :preferences, updated_at = CURRENT_TIMESTAMP
			WHERE id = :id
			RETURNING %s`, defReturnCols)

	return r.update(ctx, query, user)
}

func (r *uRepository) update(ctx context.Context, query string, user users.User) (users.User, error) {
	dUser, err := toDBUser(user)
	if err != nil {
		return users.User{}, err
	}

	rows, err := r.NamedQueryContext(ctx, query, dUser)
	if err != nil {
		return users.User{}, err
	}
	defer rows.Close()

	dUser = dbUser{}
	if rows.Next() {
		if err := rows.StructScan(&dUser); err != nil {
			return users.User{}, err
		}

		return fromDBUser(dUser), nil
	}

	return users.User{}, errors.New("user not found")
}

func (r *uRepository) Delete(ctx context.Context, id string) error {
	q := "DELETE FROM users WHERE id = $1;"

	result, err := r.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Join(errors.New("could not delete user"), err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

type dbUser struct {
	ID          string           `db:"id"`
	Username    string           `db:"username"`
	DisplayName string           `db:"display_name"`
	Bio         string           `db:"bio"`
	PictureURL  string           `db:"picture_url"`
	Email       string           `db:"email"`
	Password    string           `db:"password"`
	Preferences pgtype.TextArray `db:"preferences"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at"`
}

func toDBUser(user users.User) (dbUser, error) {
	var preferences pgtype.TextArray
	if err := preferences.Set(user.Preferences); err != nil {
		return dbUser{}, err
	}

	return dbUser{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		PictureURL:  user.PictureURL,
		Email:       user.Email,
		Password:    user.Password,
		Preferences: preferences,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func fromDBUser(dbUser dbUser) users.User {
	preferences := make([]string, len(dbUser.Preferences.Elements))
	for i, element := range dbUser.Preferences.Elements {
		preferences[i] = element.String
	}

	return users.User{
		ID:          dbUser.ID,
		Username:    dbUser.Username,
		DisplayName: dbUser.DisplayName,
		Bio:         dbUser.Bio,
		PictureURL:  dbUser.PictureURL,
		Email:       dbUser.Email,
		Password:    dbUser.Password,
		Preferences: preferences,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
	}
}
