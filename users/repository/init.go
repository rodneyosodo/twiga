// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
package repository

import (
	_ "github.com/jackc/pgx/v5/stdlib" // required for SQL access
	migrate "github.com/rubenv/sql-migrate"
)

func Migration() *migrate.MemoryMigrationSource {
	return &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "users_01",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS users (
						id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
						username VARCHAR(1024) UNIQUE NOT NULL CHECK (username <> ''),
						display_name VARCHAR(1024) NOT NULL CHECK (display_name <> ''),
						bio TEXT NOT NULL,
						picture_url VARCHAR NOT NULL,
						email VARCHAR(1024) UNIQUE NOT NULL CHECK (email <> ''),
						password TEXT NOT NULL CHECK (password <> ''),
						preferences TEXT[] NOT NULL,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
					)`,
					`CREATE TABLE IF NOT EXISTS preferences (
						id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
						user_id UUID UNIQUE NOT NULL,
						email_enabled BOOLEAN NOT NULL,
						push_enabled BOOLEAN NOT NULL,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
					)`,
					`CREATE TABLE IF NOT EXISTS followers (
						id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
						follower_id UUID NOT NULL,
						followee_id UUID NOT NULL,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						UNIQUE (follower_id, followee_id),
						CHECK (follower_id <> followee_id),
						FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
						FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE
					)`,
					`CREATE TABLE IF NOT EXISTS feeds (
						id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
						user_id UUID NOT NULL,
						post_id UUID NOT NULL,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						UNIQUE (user_id, post_id),
						CHECK (user_id <> post_id),
						FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
					)`,
					`CREATE INDEX idx_users_username ON users(username);`,
					`CREATE INDEX idx_users_email ON users(email);`,
					`CREATE INDEX idx_preferences_user_id ON preferences(user_id);`,
					`CREATE INDEX idx_preferences_email_enabled ON preferences(email_enabled);`,
					`CREATE INDEX idx_preferences_push_enabled ON preferences(push_enabled);`,
					`CREATE INDEX idx_followers_follower_id ON followers(follower_id);`,
					`CREATE INDEX idx_followers_followee_id ON followers(followee_id);`,
					`CREATE INDEX idx_feeds_user_id ON feeds(user_id);`,
					`CREATE INDEX idx_feeds_post_id ON feeds(post_id);`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS feeds CASCADE`,
					`DROP TABLE IF EXISTS followers CASCADE`,
					`DROP TABLE IF EXISTS preferences CASCADE`,
					`DROP TABLE IF EXISTS users CASCADE`,
				},
			},
		},
	}
}
