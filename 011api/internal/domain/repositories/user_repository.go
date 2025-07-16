package repositories

import (
	"context"

	"github.com/golang/011api/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (*entities.User, error)
	GetById(ctx context.Context, id int) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Update(ctx context.Context, id int, user *entities.UserInput) (*entities.UserInput, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
	Count(ctx context.Context) (int, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
}