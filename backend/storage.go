package backend

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct{ SQL *sql.DB }

func Open(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &DB{SQL: db}, nil
}

func (d *DB) Close() error { return d.SQL.Close() }

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS repertoire (
      id       INTEGER PRIMARY KEY AUTOINCREMENT,
      name     TEXT NOT NULL UNIQUE,
      color    TEXT NOT NULL CHECK (color IN ('white','black')),
      elo      INTEGER NOT NULL DEFAULT 1200,
      coverage REAL NOT NULL DEFAULT 0.0
    );
    CREATE TABLE IF NOT EXISTS nodes (
      fen         TEXT NOT NULL,
      rep_id      INTEGER NOT NULL,
      sr_index    INTEGER NOT NULL DEFAULT 0,
      due         INTEGER,
      last_review INTEGER,
      PRIMARY KEY (fen, rep_id),
      FOREIGN KEY (rep_id) REFERENCES repertoire(id) ON DELETE CASCADE
    );
    CREATE TABLE IF NOT EXISTS stats (
      fen        TEXT PRIMARY KEY,
      games      INTEGER NOT NULL DEFAULT 0,
      white_win  REAL NOT NULL DEFAULT 0.0,
      black_win  REAL NOT NULL DEFAULT 0.0,
      draw       REAL NOT NULL DEFAULT 0.0,
      moves      TEXT NOT NULL DEFAULT '[]'
    );
    CREATE TABLE IF NOT EXISTS edges (
      rep_id     INTEGER NOT NULL,
      parent_fen TEXT NOT NULL,
      child_fen  TEXT NOT NULL,
      move       TEXT NOT NULL,
      PRIMARY KEY (rep_id, parent_fen, child_fen),
      FOREIGN KEY (rep_id) REFERENCES repertoire(id) ON DELETE CASCADE
    );
    `)
	return err
}
