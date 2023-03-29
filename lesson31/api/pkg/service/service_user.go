package service

import (
	"context"
	"http_service/pkg/repository"
	"http_service/pkg/user"
)

type ServiceUser struct {
	repo repository.RepositoryUser
}

func NewServiceUser(repo repository.RepositoryUser) *ServiceUser {
	return &ServiceUser{repo: repo}
}

func (s *ServiceUser) CreateUser(ctx context.Context, user *user.User) error {
	return s.repo.CreateUser(ctx, user)
}
func (s *ServiceUser) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}
func (s *ServiceUser) GetUserFriend(ctx context.Context, id int64) ([]*user.Friend, error) {
	return s.repo.GetUserFriend(ctx, id)
}
func (s *ServiceUser) GetUser(ctx context.Context, id int64) (*user.User, error) {
	return s.repo.GetUser(ctx, id)
}
func (s *ServiceUser) UpdateUser(ctx context.Context, user *user.User) error {
	return s.repo.UpdateUser(ctx, user)
}
