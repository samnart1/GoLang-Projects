package repositories

import (
	"context"

	"github.com/golang/011api/internal/domain/entities"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entities.Comment, postID, userID int) (*entities.Comment, error)
	GetByID(ctx context.Context, id int) (*entities.Comment, error)
	Update(ctx context.Context, id int, comment *entities.CommentInput) (*entities.Comment, error)
	Delete(ctx context.Context, id int) error
	ListByPost(ctx context.Context, postID, limit, offset int) ([]*entities.Comment, error)
	ListByUser(ctx context.Context, userID, limit, offset int) ([]*entities.Comment, error)
	Count(ctx context.Context) (int, error)
	CountByPost(ctx context.Context, postID int) (int, error)
	CountByUser(ctx context.Context, userID int) (int, error)
}