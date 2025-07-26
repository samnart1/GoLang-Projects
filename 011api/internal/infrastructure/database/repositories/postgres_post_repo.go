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

type PostgresPostRepository struct {
	db *database.DB
}


func (r *PostgresPostRepository) Count(ctx context.Context) (int, error) {
	// panic("unimplemented")
	query := `
		SELECT COUNT(*) FROM posts
	`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}

	return count, nil
}


func (r *PostgresPostRepository) CountByAuthor(ctx context.Context, AuthorID int) (int, error) {
	// panic("unimplemented")
	query := `
		SELECT COUNT(*) FROM posts WHERE author_id = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, AuthorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts by author: %w", err)
	}

	return count, nil
}


func (r *PostgresPostRepository) CountPublished(ctx context.Context) (int, error) {
	// panic("unimplemented")
	query := `
		SELECT COUNT(*) FROM posts WHERE published = true
	`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count published posts: %w", err)
	}

	return count, nil
}


func (r *PostgresPostRepository) Create(ctx context.Context, postInput *entities.PostInput, authorID int) (*entities.Post, error) {
	// panic("unimplemented")
	query := `
		INSERT INTO posts (title, content, slug, author_id, published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, content, slug, author_id, published, published_at, created_at, updated_at
	`
	now := time.Now()
	post := &entities.Post{}

	err := r.db.QueryRowContext(ctx, query,
		postInput.Title,
		postInput.Content,
		postInput.Slug,
		authorID,
		false,
		now,
		now,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Slug,
		&post,authorID,
		&post.Published,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}


func (r *PostgresPostRepository) Delete(ctx context.Context, id int) error {
	// panic("unimplemented")
	query := `DELETE FROM posts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}


func (r *PostgresPostRepository) GetByID(ctx context.Context, id int) (*entities.Post, error) {
	// panic("unimplemented")
	query := `
		SELECT p.id, p.title, p.content, p.slug, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
			u.id, u.email, u.username, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM posts p
		LEFT JOIN users u ON P.author_id = u.id
		WHERE p.id = $1
	`

	post := &entities.Post{}
	var author entities.User
	var authorID sql.NullInt64
	var authorEmail, authorUsername sql.NullString
	var authorFirstName, authorLastName sql.NullString
	var authorIsActive sql.NullBool
	var authorCreatedAt, authorUpdatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Slug,
		&post.AuthorID,
		&post.Published,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
		&authorID,
		&authorEmail,
		&authorUsername,
		&authorFirstName,
		&authorLastName,
		&authorIsActive,
		&authorCreatedAt,
		&authorUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if authorID.Valid {
		author.ID = int(authorID.Int64)
		author.Email = authorEmail.String
		author.Username = authorUsername.String
		if authorFirstName.Valid {
			author.FirstName = &authorFirstName.String
		}
		if authorLastName.Valid {
			author.LastName = &authorLastName.String
		}
		author.IsActive = authorIsActive.Bool
		author.CreatedAt = authorCreatedAt.Time
		author.UpdatedAt = authorUpdatedAt.Time
		post.Author = &author
	}

	return post, nil
}


func (r *PostgresPostRepository) GetBySlug(ctx context.Context, slug string) (*entities.Post, error) {
	// panic("unimplemented")
	query := `
		SELECT p.id, p.title, p.content, p.slug, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
			u.id, u.email, u.username, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.slug = $1
	`

	post := &entities.Post{}
	var author entities.User
	var authorID sql.NullInt64
	var authorEmail, authorUsername sql.NullString
	var authorFirstName, authorLastName sql.NullString
	var authorIsActive sql.NullBool
	var authorCreatedAt, authorUpdatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Slug,
		&post.AuthorID,
		&post.Published,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
		&authorID,
		&authorEmail,
		&authorUsername,
		&authorFirstName,
		&authorLastName,
		&authorIsActive,
		&authorCreatedAt,
		&authorUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if authorID.Valid {
		author.ID = int(authorID.Int64)
		author.Email = authorEmail.String
		author.Username = authorUsername.String
		if authorFirstName.Valid {
			author.FirstName = &authorFirstName.String
		}
		if authorLastName.Valid {
			author.LastName = &authorLastName.String
		}
		author.IsActive = authorIsActive.Bool
		author.CreatedAt = authorCreatedAt.Time
		author.UpdatedAt = authorUpdatedAt.Time
		post.Author = &author
	}

	return post, nil
}


func (r *PostgresPostRepository) List(ctx context.Context, limit int, offset int) ([]*entities.Post, error) {
	// panic("unimplemented")
	query := `
		SELECT p.id, p.title, p.content, p.slug, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
			u.id, u.email, u.username, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	return r.scanPosts(ctx, query, limit, offset)
}


func (r *PostgresPostRepository) ListByAuthor(ctx context.Context, AuthorID int, limit int, offset int) ([]*entities.Post, error) {
	// panic("unimplemented")
	query := `
		SELECT p.id, p.title, p.content, p.slug, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
			u.id, u.email, u.username, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.author_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	return  r.scanPosts(ctx, query, limit, offset)
}


func (r *PostgresPostRepository) ListPublished(ctx context.Context, limit int, offset int) ([]*entities.Post, error) {
	// panic("unimplemented")
	query := `
		SELECT p.id, p.title, p.content, p.slug, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
			u.id, u.email, u.username, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.published = true
		ORDER BY p.published_at DESC
		LIMIT $1 OFFSET $2
	`

	return r.scanPosts(ctx, query, limit, offset)
}

func (r *PostgresPostRepository) scanPosts(ctx context.Context, query string, args ...interface{}) ([]*entities.Post, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return  nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		post := &entities.Post{}
		var author entities.User
		var authorID sql.NullInt64
		var authorEmail, authorUsername sql.NullString
		var authorFirstName, authorLastName sql.NullString
		var authorIsActive sql.NullBool
		var authorCreatedAt, authorUpdatedAt sql.NullTime

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Slug,
			&post.AuthorID,
			&post.Published,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
			&authorID,
			&authorEmail,
			&authorUsername,
			&authorFirstName,
			&authorLastName,
			&authorIsActive,
			&authorCreatedAt,
			&authorUpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		if authorID.Valid {
			author.ID = int(authorID.Int64)
			author.Email = authorEmail.String
			author.Username = authorUsername.String
			if authorFirstName.Valid {
				author.FirstName = &authorFirstName.String
			}
			if authorLastName.Valid {
				author.LastName = &authorLastName.String
			}
			author.IsActive = authorIsActive.Bool
			author.CreatedAt = authorCreatedAt.Time
			author.UpdatedAt = authorUpdatedAt.Time
			post.Author = &author
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}


func (r *PostgresPostRepository) Publish(ctx context.Context, id int) error {
	// panic("unimplemented")
	query := `
		UPDATE posts
		SET published = true, published_at = $2, updated_at = $3
		WHERE id = $1
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, id, now, now)
	if err != nil {
		return fmt.Errorf("failed to publish post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return  fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}


func (r *PostgresPostRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	// panic("unimplemented")
	query := `
		SELECT EXISTS(SELECT 1 FROM posts WHERE slug = $1)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check slug exists: %w", err)
	}

	return exists, nil
}


func (r *PostgresPostRepository) UnPublish(ctx context.Context, id int) error {
	// panic("unimplemented")
	query := `
		UPDATE posts
		SET published = false, published_at = NULL, updated_at = $2
		WHERE id = $1
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("failed to unpublish post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}


func (r *PostgresPostRepository) Update(ctx context.Context, id int, postInput *entities.PostInput) (*entities.Post, error) {
	// panic("unimplemented")
	query := `
		UPDATE posts
		SET title = $2, content = $3, slug = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, title, content, slug, author_id, published, published_at, created_at, updated_at
	`

	post := &entities.Post{}
	err := r.db.QueryRowContext(ctx, query,
		id,
		postInput.Title,
		postInput.Content,
		postInput.Slug,
		time.Now(),
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Slug,
		&post.AuthorID,
		&post.Published,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func NewPostgresPostRepository(db *database.DB) repositories.PostRepository {
	return &PostgresPostRepository{db: db}
}
