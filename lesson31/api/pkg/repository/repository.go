package repository

import (
	"context"
	"http_service/pkg/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryUser interface {
	CreateUser(context.Context, *user.User) error
	UpdateUser(context.Context, *user.User) error
	DeleteUser(context.Context, int64) error
	GetUserFriend(context.Context, int64) ([]*user.Friend, error)
	GetUser(context.Context, int64) (*user.User, error)
}

type Repository struct {
	RepositoryUser
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		RepositoryUser: NewRepositoryUser(client),
	}
}
