package repository

import (
	"context"
	"database/sql"
	"fmt"
	"translator/domain"
)

type flapRepository struct {
	database *sql.DB
}

func NewFlapRepository(db *sql.DB) domain.FlapRepository {
	return &flapRepository{
		database: db,
	}
}

func (fr *flapRepository) Add(ctx context.Context, record *domain.DefaultRecord) error {
	tx, err := fr.database.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	stmt := `
	INSERT INTO flap (timestamp, position, value) VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(stmt, record.Timestamp, record.Position, record.Value)
	if err != nil {
		return fmt.Errorf("failed to insert flap record: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit flap insertion: %v", err)
	}
	return nil
}
