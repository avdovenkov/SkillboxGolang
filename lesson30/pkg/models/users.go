package models

import (
	"errors"
	"fmt"
)

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []*User `json:"friends"`
}

type UserInRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type FriendRequest struct {
	Source_id int `json:"source_id"`
	Target_id int `json:"target_id"`
}

func (u *User) MakeFrienship(target *User) (*User, error) {
	if u == target {
		return nil, errors.New(fmt.Sprintf("can't friend with himself: %v", u))
	}
	u.Friends = append(u.Friends, target)
	return u, nil
}
