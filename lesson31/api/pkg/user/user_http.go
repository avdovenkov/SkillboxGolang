package user

import "net/http"

type UserRequest struct {
	*User
}
type FriendRequest struct {
	Source_id int64 `json:"source_id"`
	Target_id int64 `json:"target_id"`
}
type UserAgeUpdate struct {
	NewAge int64 `json:"new_age"`
}

func (u *UserRequest) Bind(r *http.Request) error {
	return nil
}
func (u *UserAgeUpdate) Bind(r *http.Request) error {
	return nil
}
func (u *FriendRequest) Bind(r *http.Request) error {
	return nil
}

