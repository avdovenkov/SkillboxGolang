package storage

import (
	"errors"
	"fmt"
	"httpServer/pkg/models"
)

type MemStorage map[int]*models.User

var cnt int = 1

func NewStorage() *MemStorage {
	var storage MemStorage = make(map[int]*models.User)
	return &storage
}

func (s *MemStorage) Get(key int) *models.User {
	user, ok := (*s)[key]
	if !ok {
		return nil
	}
	return user
}

func (s *MemStorage) Save(u *models.User) error {
	for {
		_, ok := (*s)[cnt]
		if ok {
			cnt++
			continue
		}
		break
	}
	u.ID = cnt
	(*s)[cnt] = u
	return nil
}

func (s *MemStorage) AllUsers() []models.User {
	var users []models.User
	for _, u := range *s {
		users = append(users, *u)
	}
	return users
}

func (s *MemStorage) Delete(id int) error {
	if user := s.Get(id); user == nil {
		return errors.New(fmt.Sprintf("user %d not found", id))
	}
	delete(*s, id)
	s.DeleteFriend(id)
	return nil
}

func (s *MemStorage) Update(id int, u *models.User) (*models.User, error) {
	_, ok := (*s)[id]
	if !ok {
		return &models.User{}, errors.New(fmt.Sprintf("user %d not found", id))
	}
	(*s)[id] = u
	return (*s)[id], nil
}

func (s *MemStorage) DeleteFriend(id int) {
	for _, u := range *s {
		for i, friend := range u.Friends {
			if friend.ID == id {
				u.Friends = append(u.Friends[:i], u.Friends[i+1:]...)
			}
		}
	}
}
