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

	// CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	_, err = tx.Exec(`
	DROP TABLE IF EXISTS users;
	DROP TABLE IF EXISTS user_roles;
	CREATE TABLE IF NOT EXISTS user_roles (
		role_id INTEGER PRIMARY KEY,
		title TEXT
	);

	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		role_id INTEGER REFERENCES user_roles,
		email TEXT UNIQUE,
		passwd TEXT,
		creation_date DATE
	);

	CREATE INDEX IF NOT EXISTS idx_users_email
		ON public.users USING HASH
		(email COLLATE pg_catalog."default")
		TABLESPACE pg_default;

	INSERT INTO user_roles (role_id, title) VALUES (0, 'admin');
	INSERT INTO user_roles (role_id, title) VALUES (1, 'employee');
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
