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
				Id: "notifications_01",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS notifications (
						id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
						user_id UUID NOT NULL,
						category VARCHAR(256) NOT NULL CHECK (category <> ''),
						content TEXT NOT NULL,
						is_read BOOLEAN NOT NULL DEFAULT FALSE,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
					)`,
					`CREATE INDEX idx_notifications_user_id ON notifications(user_id);`,
					`CREATE INDEX idx_notifications_category ON notifications(category);`,
					`CREATE INDEX idx_notifications_is_read ON notifications(is_read);`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS notifications CASCADE`,
				},
			},
		},
	}
}
