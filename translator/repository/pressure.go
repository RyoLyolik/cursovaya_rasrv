package repository

import (
	"context"
	"database/sql"
	"fmt"
	"translator/domain"
)

type pressureRepository struct {
	database *sql.DB
}

func NewPressureRepository(db *sql.DB) domain.PressureRepository {
	return &pressureRepository{
		database: db,
	}
}

func (pr *pressureRepository) Add(ctx context.Context, record *domain.DefaultRecord) error {
	tx, err := pr.database.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	stmt := `
	INSERT INTO pressure (timestamp, position, value) VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(stmt, record.Timestamp, record.Position, record.Value)
	if err != nil {
		return fmt.Errorf("failed to insert pressure record: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit pressure insertion: %v", err)
	}
	return nil
}
