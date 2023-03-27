package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"httpServer/pkg/models"
	"httpServer/pkg/server"
	"io"
	"log"
	"net/http"
	"strconv"
)

func AddUserToServer(name string, age int) ([]byte, error) {
	user := models.UserInRequest{Name: name, Age: age}
	log.Printf("[CLIENT] adding new user: %v", user)

	json_data, err := json.Marshal(user)
	if err != nil {
		explain := fmt.Sprintf("can't convert json data: %v", err)
		return nil, errors.New(explain)
	}

	resp, err := http.Post("http://127.0.0.1:8080/users/create", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func GetUserFromServer(id int) {
	log.Printf("[CLIENT] requesting user %d\n", id)
	resp, err := http.Get("http://127.0.0.1:8080/users/" + strconv.Itoa(id))
	if err != nil {
		log.Printf("[CLIENT] can't get user %d: %v", id, err)
	}
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}

func GetAllUserFromServer() {
	log.Println("[CLIENT] requesting all users")
	resp, err := http.Get("http://127.0.0.1:8080/users/")
	if err != nil {
		log.Printf("[CLIENT] can't get users: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	log.Println(string(body))
}

func DeleteUserFromServer(id int) {
	log.Printf("[CLIENT] deleting user %d\n", id)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://127.0.0.1:8080/users/"+strconv.Itoa(id), nil)
	if err != nil {
		log.Printf("[CLIENT] can't get response: %v", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[CLIENT] can't get response: %v", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}

func UpdateUserAtServer(id int, u models.User) {
	log.Printf("[CLIENT] updating user %d\n", id)
	jsonData, err := json.Marshal(u)
	if err != nil {
		log.Printf("[CLIENT] can't parse user id %d: %v", id, err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:8080/users/"+strconv.Itoa(id), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[CLIENT] can't get response: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[CLIENT] can't get response: %v", err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}

func GetUserFromResponse(resp []byte) (models.User, error) {
	var jsonData server.UserResponse
	err := json.Unmarshal(resp, &jsonData)
	if err != nil {
		log.Printf("[CLIENT] can't parse response: %v\n", err)
		return models.User{}, err
	}
	return jsonData.Response, nil
}

func MakeFriendsOnServer(source, target models.User) ([]byte, error) {
	log.Printf("[CLIENT] adding friend %d for user %d", target.ID, source.ID)

	request := models.FriendRequest{
		Source_id: source.ID,
		Target_id: target.ID,
	}
	json_data, err := json.Marshal(request)
	if err != nil {
		explain := fmt.Sprintf("can't convert json data: %v", err)
		return nil, errors.New(explain)
	}

	resp, err := http.Post("http://127.0.0.1:8080/users/make_friend", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func GetFriendsFromServer(id int) {
	log.Printf("[CLIENT] requesting friends for %d user\n", id)
	resp, err := http.Get("http://127.0.0.1:8080/friends/" + strconv.Itoa(id))
	if err != nil {
		log.Printf("[CLIENT] can't get user %d: %v", id, err)
	}
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}
