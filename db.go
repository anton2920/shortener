package main

import (
	"database/sql"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"

	"github.com/anton2920/gofa/l10n"
)

type Database struct {
	*sql.DB
	l10n.Language
}

var DB *sql.DB

func OpenDB(dsn string) error {
	var err error

	DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open sqlite3 DB: %v", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping sqlite3 DB: %v", err)
	}

	if _, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		flags INT NOT NULL DEFAULT 0,

		email TEXT NOT NULL UNIQUE,
		password BYTEA NOT NULL,

		created_at BIGINT NOT NULL
	)`); err != nil {
		return fmt.Errorf("failed to create table for users: %v", err)
	}

	return nil
}

func UniqueViolation(err error) bool {
	if sqlErr, ok := err.(sqlite3.Error); ok {
		return sqlErr.ExtendedCode == 2067
	}
	return false
}
