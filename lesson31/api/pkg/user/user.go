package user

import (
	"fmt"
	"net/http"
)

type User struct {
	Id      int64   `json:"id,omitempty" bson:"id"`
	Name    string  `json:"name"`
	Age     int64   `json:"age"`
	Friends []int64 `json:"friends"`
}
type Friend struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func (u *User) NewFiend(user *User) error {
	for _, one_friend := range u.Friends {
		if user.Id == one_friend {
			return fmt.Errorf("пользователь с id %d уже является другом пользователя %d", user.Id, u.Id)
		}
	}
	u.Friends = append(u.Friends, user.Id)
	user.Friends = append(user.Friends, u.Id)
	return nil
}
func (u *User) ToString() string {
	return fmt.Sprintf("name:%s,age:%d", u.Name, u.Age)
}
func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (u *Friend) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (u *User) Bind(r *http.Request) error {
	return nil
}
