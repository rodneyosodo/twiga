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
	"time"

	"github.com/jackc/pgtype"
	"github.com/rodneyosodo/twiga/internal/postgres"
	"github.com/rodneyosodo/twiga/notifications"
)

var _ notifications.Repository = (*repository)(nil)

type repository struct {
	postgres.Database
}

func NewRepository(db postgres.Database) *repository {
	return &repository{db}
}

func (r *repository) CreateNotification(ctx context.Context, notification notifications.Notification) (notifications.Notification, error) {
	query := `INSERT INTO notifications (user_id, category, content, is_read)
			VALUES (:user_id, :category, :content, :is_read)
			RETURNING *`
	dNotification, err := toDBNotification(notification)
	if err != nil {
		return notifications.Notification{}, err
	}

	rows, err := r.NamedQueryContext(ctx, query, dNotification)
	if err != nil {
		return notifications.Notification{}, err
	}
	defer rows.Close()

	dNotification = dbNotification{}
	if rows.Next() {
		if err := rows.StructScan(&dNotification); err != nil {
			return notifications.Notification{}, err
		}
	}

	return dNotification.toNotification(), nil
}

func (r *repository) RetrieveNotification(ctx context.Context, id string) (notifications.Notification, error) {
	query := `SELECT * FROM notifications WHERE id = :id`
	dNotification := dbNotification{
		ID: id,
	}

	rows, err := r.NamedQueryContext(ctx, query, dNotification)
	if err != nil {
		return notifications.Notification{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var dNotification dbNotification
		if err := rows.StructScan(&dNotification); err != nil {
			return notifications.Notification{}, err
		}

		return dNotification.toNotification(), nil
	}

	return notifications.Notification{}, errors.New("notification not found")
}

func (r *repository) RetrieveAllNotifications(ctx context.Context, page notifications.Page) (npage notifications.NotificationsPage, err error) {
	filter := ""
	filters := []string{}
	if page.UserID != "" {
		filters = append(filters, fmt.Sprintf("user_id = '%s'", page.UserID))
	}
	if page.Category.String() != "" {
		filters = append(filters, fmt.Sprintf("category = '%s'", page.Category.String()))
	}
	if page.IsRead != nil {
		filters = append(filters, fmt.Sprintf("is_read = %t", *page.IsRead))
	}
	if page.IDs != nil {
		filters = append(filters, fmt.Sprintf("user_id IN ('%s')", strings.Join(page.IDs, "','")))
	}
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " AND "))
	}

	query := fmt.Sprintf(`SELECT * FROM notifications %s ORDER BY created_at DESC LIMIT :limit OFFSET :offset`, filter)

	dPage := notifications.NotificationsPage{
		Page: page,
	}

	rows, err := r.NamedQueryContext(ctx, query, dPage)
	if err != nil {
		return notifications.NotificationsPage{}, err
	}
	defer rows.Close()

	items := make([]notifications.Notification, 0)
	for rows.Next() {
		var dNotification dbNotification
		if err := rows.StructScan(&dNotification); err != nil {
			return notifications.NotificationsPage{}, err
		}

		items = append(items, dNotification.toNotification())
	}

	totalQuery := fmt.Sprintf(`SELECT COUNT(*) FROM notifications %s`, filter)

	total, err := postgres.Total(ctx, r.Database, totalQuery, dPage)
	if err != nil {
		return notifications.NotificationsPage{}, err
	}

	return notifications.NotificationsPage{
		Page: notifications.Page{
			Limit:  page.Limit,
			Offset: page.Offset,
			Total:  total,
		},
		Notifications: items,
	}, nil
}

func (r *repository) ReadNotification(ctx context.Context, userID, id string) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = :id AND user_id = :user_id`
	dNotification := dbNotification{
		ID:     id,
		UserID: userID,
	}

	result, err := r.NamedExecContext(ctx, query, dNotification)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("could not read notification")
	}

	return nil
}

func (r *repository) ReadAllNotifications(ctx context.Context, page notifications.Page) error {
	filter := ""
	filters := []string{}
	if page.UserID != "" {
		filters = append(filters, fmt.Sprintf("user_id = '%s'", page.UserID))
	}
	if page.Category.String() != "" {
		filters = append(filters, fmt.Sprintf("category = '%s'", page.Category.String()))
	}
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " AND "))
	}

	query := fmt.Sprintf(`UPDATE notifications SET is_read = TRUE %s`, filter)
	dPage := notifications.NotificationsPage{
		Page: page,
	}

	result, err := r.NamedExecContext(ctx, query, dPage)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("could not read notifications")
	}

	return nil
}

func (r *repository) DeleteNotification(ctx context.Context, id string) error {
	query := `DELETE FROM notifications WHERE id = :id`
	dNotification := dbNotification{
		ID: id,
	}

	result, err := r.NamedExecContext(ctx, query, dNotification)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("could not delete notification")
	}

	return nil
}

type dbNotification struct {
	ID        string      `db:"id"`
	UserID    string      `db:"user_id"`
	Category  string      `db:"category"`
	Content   string      `db:"content"`
	IsRead    pgtype.Bool `db:"is_read"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

func (n dbNotification) toNotification() notifications.Notification {
	return notifications.Notification{
		ID:        n.ID,
		UserID:    n.UserID,
		Category:  notifications.ToCategory(n.Category),
		Content:   n.Content,
		IsRead:    n.IsRead.Bool,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

func toDBNotification(n notifications.Notification) (dbNotification, error) {
	var isRead pgtype.Bool
	if err := isRead.Set(n.IsRead); err != nil {
		return dbNotification{}, err
	}

	return dbNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		Category:  n.Category.String(),
		Content:   n.Content,
		IsRead:    isRead,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}, nil
}
