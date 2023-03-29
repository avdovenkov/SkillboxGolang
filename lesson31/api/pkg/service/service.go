package service

import (
	"context"
	"http_service/pkg/repository"
	"http_service/pkg/user"
)

type ServiceUsers interface {
	CreateUser(context.Context, *user.User) error
	UpdateUser(context.Context, *user.User) error
	DeleteUser(context.Context, int64) error
	GetUserFriend(context.Context, int64) ([]*user.Friend, error)
	GetUser(context.Context, int64) (*user.User, error)
}

type Service struct {
	ServiceUsers
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		ServiceUsers: NewServiceUser(repo)}
}
