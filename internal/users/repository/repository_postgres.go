package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Postgres is a repository implementation for Postgres.
type Postgres struct {
	db *sqlx.DB
}

// NewPostgres creates a new Postgres repository.
func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

// Get returns a user by id.
func (p *Postgres) Get(ctx context.Context, id string) (*User, error) {
	var user User
	if err := p.db.GetContext(
		ctx,
		&user,
		`SELECT id, first_name, last_name, nickname, password, email,
		country, created_at, updated_at FROM users WHERE id =$1`,
		id,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("could not get user: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("could not get user: %w", err)
	}
	return &user, nil
}

func (p *Postgres) GetAll(ctx context.Context, cursor string, limit int) ([]*User, error) {
	var users []*User
	if cursor == "" {
		if err := p.db.SelectContext(
			ctx,
			&users,
			`SELECT id, first_name, last_name, nickname, password, email, country, 
			created_at, updated_at FROM users ORDER BY id ASC LIMIT $1`,
			limit,
		); err != nil {
			return nil, fmt.Errorf("could not get users: %w", err)
		}
		return users, nil
	}

	if err := p.db.SelectContext(
		ctx,
		&users,
		`SELECT id, first_name, last_name, nickname, password, email, country,  
		created_at, updated_at FROM users WHERE id > $1 ORDER BY id ASC LIMIT $2`,
		cursor,
		limit,
	); err != nil {
		return nil, fmt.Errorf("could not get users: %w", err)
	}
	return users, nil
}

// GetByCountry returns a list of users by country.
func (p *Postgres) GetByCountry(ctx context.Context, country string, cursor string, limit int) ([]*User, error) {
	var users []*User
	if cursor == "" {
		if err := p.db.SelectContext(
			ctx,
			&users,
			`SELECT id, first_name, last_name, nickname, password, email, country,
			created_at, updated_at FROM users WHERE country = $1 ORDER BY id ASC LIMIT $2`,
			country,
			limit,
		); err != nil {
			return nil, fmt.Errorf("could not get users: %w", err)
		}
		return users, nil
	}

	if err := p.db.SelectContext(
		ctx,
		&users,
		`SELECT id, first_name, last_name, nickname, password, email, country, created_at, 
		updated_at FROM users WHERE country= $1 AND id > $2 ORDER BY id ASC LIMIT $3`,
		country,
		cursor,
		limit,
	); err != nil {
		return nil, fmt.Errorf("could not get users: %w", err)
	}
	return users, nil
}

// Insert inserts a new user.
func (p *Postgres) Insert(ctx context.Context, user *User) error {
	if _, err := p.db.NamedExecContext(
		ctx,
		`INSERT INTO users (id, first_name, last_name, nickname, password, email, country, created_at, updated_at) 
		VALUES (:id, :first_name, :last_name, :nickname, :password, :email, :country, :created_at, :updated_at)`,
		user,
	); err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" { // unique_violation: https://www.postgresql.org/docs/8.2/errcodes-appendix.html
				return fmt.Errorf("could not insert user: %w", ErrDuplicateEmail)
			}
		}
		return fmt.Errorf("could not insert user: %w", err)
	}
	return nil
}

// Update updates a user by id.
func (p *Postgres) Update(ctx context.Context, user *User) error {
	result, err := p.db.NamedExecContext(
		ctx,
		`UPDATE users SET first_name = :first_name, last_name = :last_name, nickname = :nickname, 
		password = :password, email = :email, country = :country, updated_at = :updated_at WHERE id = :id`,
		user,
	)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("could not insert user: %w", ErrDuplicateEmail)
			}
		}
		return fmt.Errorf("could not update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("could not update user: %w", ErrUserNotFound)
	}
	return nil
}

// Delete deletes a user by id.
func (p *Postgres) Delete(ctx context.Context, id string) error {
	if _, err := p.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("could not delete user: %w", ErrUserNotFound)
		}
		return fmt.Errorf("could not delete user: %w", err)
	}
	return nil
}

// CheckDatabaseHealth checks if the database is healthy by pinging it.
func (p *Postgres) CheckDatabaseHealth(ctx context.Context) error {
	if err := p.db.PingContext(ctx); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}
	return nil
}
