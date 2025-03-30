package bootstrap

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func newPostgres(log *slog.Logger, psqlConnectionInfo string) (*sql.DB, error) {
	log.Info("Connecting to postgres")
	db, err := sql.Open("postgres", psqlConnectionInfo)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return db, nil
}

func MakePostgres(log *slog.Logger, psqlConnectionInfo string) (*sql.DB, error) {
	db, err := newPostgres(log, psqlConnectionInfo)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	_, err = tx.Exec(`
	DROP TABLE IF EXISTS temperature;
	DROP TABLE IF EXISTS pressure;
	DROP TABLE IF EXISTS flap;
	CREATE TABLE IF NOT EXISTS temperature (
		record_id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP NOT NULL,
		position INTEGER NOT NULL,
		value DOUBLE PRECISION NOT NULL
	);
	CREATE TABLE IF NOT EXISTS pressure (
		record_id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP NOT NULL,
		position INTEGER NOT NULL,
		value DOUBLE PRECISION NOT NULL
	);
	CREATE TABLE IF NOT EXISTS flap (
		record_id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP NOT NULL,
		position INTEGER NOT NULL,
		value DOUBLE PRECISION NOT NULL
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return db, nil
}
