package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/011api/internal/domain/entities"
	"github.com/golang/011api/internal/domain/repositories"
	"github.com/golang/011api/internal/infrastructure/database"
)

type PostgresUserRepository struct {
	db *database.DB
}


func (r *PostgresUserRepository) Count(ctx context.Context) (int, error) {
	// panic("unimplemented")
	query := `
		SELECT COUNT(*) FROM users
	`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return  0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}


func (r *PostgresUserRepository) Create(ctx context.Context, userInput *entities.User) (*entities.User, error) {
	// panic("unimplemented")
	query := `
		INSERT INTO users (email, username, password_hash, first_name, last_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, email, username, password_hash, first_name, last_name, is_active, created_at, updated_at
	`

	now := time.Now()
	user := &entities.User{}

	err := r.db.QueryRowContext(ctx, query,
		userInput.Email,
		userInput.Username,
		userInput.PasswordHash,
		userInput.FirstName,
		userInput.LastName,
		true,
		now,
		now,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return  nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}


func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	// panic("unimplemented")
	query := `
		DELETE FROM users WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return  fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}


func (r *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	// panic("unimplemented")
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}


func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	// panic("unimplemented")
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return  nil, fmt.Errorf("user not found")
		}
		return  nil, fmt.Errorf("failed to fetch user")
	}
	
	return user, nil
}


func (r *PostgresUserRepository) GetById(ctx context.Context, id int) (*entities.User, error) {
	// panic("unimplemented")
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}


func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	// panic("unimplemented")
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}


func (r *PostgresUserRepository) List(ctx context.Context, limit int, offset int) ([]*entities.User, error) {
	// panic("unimplemented")
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return  nil, fmt.Errorf("failed to list users: %w", err)
	} 
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, rows.Err()
}


func (r *PostgresUserRepository) Update(ctx context.Context, id int, userInput *entities.UserInput) (*entities.User, error) {
	// panic("unimplemented")
	query := `
		UPDATE users
		SET email = $2, username = $3, first_name = $4, last_name = $5, updated_at = $6
		WHERE id = $1
		RETURNING id, email, username, password_hash, first_name, last_name, is_active, created_at, updated_at
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query,
		id,
		userInput.Email,
		userInput.Username,
		userInput.FirstName,
		user.LastName,
		time.Now(),
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to updated user: %w", err)
	}

	return user, nil
}


func (r *PostgresUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	// panic("unimplemented")
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}

func NewPostgresUserRepository(db *database.DB) repositories.UserRepository {
	return &PostgresUserRepository{db: db}
}
