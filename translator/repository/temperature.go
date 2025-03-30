package repository

import (
	"context"
	"database/sql"
	"fmt"
	"translator/domain"
)

type temperatureRepository struct {
	database *sql.DB
}

func NewTemperatureRepository(db *sql.DB) domain.TemperatureRepository {
	return &temperatureRepository{
		database: db,
	}
}

func (tr *temperatureRepository) Add(ctx context.Context, record *domain.DefaultRecord) error {
	tx, err := tr.database.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	stmt := `
	INSERT INTO temperature (timestamp, position, value) VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(stmt, record.Timestamp, record.Position, record.Value)
	if err != nil {
		return fmt.Errorf("failed to insert temperature record: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit temperature insertion: %v", err)
	}
	return nil
}
