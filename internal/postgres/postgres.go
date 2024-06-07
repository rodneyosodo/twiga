// Copyright (c) 2024 rodneyosodo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
package postgres

import (
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // required for SQL access
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	errConnect   = errors.New("failed to connect to postgresql server")
	errMigration = errors.New("failed to apply migrations")
)

type Config struct {
	Host        string `env:"HOST"          envDefault:"localhost"`
	Port        string `env:"PORT"          envDefault:"5432"`
	User        string `env:"USER"          envDefault:"twiga"`
	Pass        string `env:"PASS"          envDefault:"twiga"`
	Name        string `env:"NAME"          envDefault:""`
	SSLMode     string `env:"SSL_MODE"      envDefault:"disable"`
	SSLCert     string `env:"SSL_CERT"      envDefault:""`
	SSLKey      string `env:"SSL_KEY"       envDefault:""`
	SSLRootCert string `env:"SSL_ROOT_CERT" envDefault:""`
}

func Setup(cfg Config, migrations migrate.MemoryMigrationSource) (*sqlx.DB, error) {
	db, err := Connect(cfg)
	if err != nil {
		return nil, err
	}

	_, err = migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		return nil, errors.Join(errMigration, err)
	}

	return db, nil
}

func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, errors.Join(errConnect, err)
	}

	return db, nil
}
