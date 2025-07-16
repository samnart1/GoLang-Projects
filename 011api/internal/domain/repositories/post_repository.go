package repositories

import (
	"context"

	"github.com/golang/011api/internal/domain/entities"
)

type PostRepository interface {
	Create(ctx context.Context, post *entities.PostInput, authorID int) (*entities.Post, error)
	GetByID(ctx context.Context, id int) (*entities.Post, error)
	GetBySlug(ctx context.Context, slug string) (*entities.Post, error)
	Update(ctx context.Context, id int, post *entities.PostInput) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) (*[]entities.Post, error)
	ListByAuthor(ctx context.Context, AuthorID, limit, offset int) (*[]entities.Post, error)
	ListPublished(ctx context.Context, limit, offset int) (*[]entities.Post, error)
	Count(ctx context.Context) (int, error)
	CountByAuthor(ctx context.Context, AuthorID int) (int, error)
	CountPublished(ctx context.Context) (int, error)
	Publish(ctx context.Context, id int) error
	UnPublish(ctx context.Context, id int) error
	SlugExists(ctx context.Context, slug string) (bool, error)
}