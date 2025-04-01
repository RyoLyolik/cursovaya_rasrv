package repository

import (
	"context"
	"database/sql"
	"fmt"
	"httpserver/domain"
	"log/slog"

	_ "github.com/lib/pq"
)

type userRepository struct {
	database *sql.DB //
	logger   *slog.Logger
}

func NewUserRepository(db *sql.DB, log *slog.Logger) domain.UserRepository {
	return &userRepository{
		database: db,
		logger:   log,
	}
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT users.user_id, users.role_id, users.email, users.passwd, users.creation_date, user_roles.title FROM users
		JOIN user_roles ON user_roles.role_id = users.role_id
		WHERE users.email = $1;
	`
	row := ur.database.QueryRow(query, email)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Role.ID,
		&user.Email,
		&user.Password,
		&user.CreationDate,
		&user.Role.Name,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
	SELECT users.user_id, users.role_id, users.email, users.passwd, users.creation_date, user_roles.title FROM users
		JOIN user_roles ON user_roles.role_id = users.role_id
		WHERE users.user_id = $1;
	`
	row := ur.database.QueryRow(query, id)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Role.ID,
		&user.Email,
		&user.Password,
		&user.CreationDate,
		&user.Role.Name,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	tx, err := ur.database.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	stmt := `
	INSERT INTO users (role_id, email, passwd, creation_date) VALUES ($1, $2, $3, $4);
	`
	_, err = tx.Exec(stmt, user.Role.ID, user.Email, user.Password, user.CreationDate)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit user insertion: %v", err)
	}
	return nil
}

func (ur *userRepository) List(ctx context.Context) ([]*domain.User, error) {
	query := `
	SELECT users.user_id, users.role_id, users.email, users.passwd, users.creation_date, user_roles.title
	FROM users JOIN user_roles ON user_roles.role_id=users.role_id;
	`
	rows, err := ur.database.Query(query)
	if err != nil {
		return nil, err
	}
	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Role.ID,
			&user.Email,
			&user.Password,
			&user.CreationDate,
			&user.Role.Name,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
