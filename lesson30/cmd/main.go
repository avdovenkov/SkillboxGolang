package main

import (
	"httpServer/pkg/client"
	"httpServer/pkg/server"
	"log"
	"net/http"
	"sync"
)

func main() {
	s := server.CreateNewServer()
	s.MountHandlers()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := http.ListenAndServe(":8080", s.Router)
		if err != nil {
			log.Fatalf("[SERVER] can't start server: %v", err)
		}
	}()

	clientRequests()

	wg.Wait()
}

func clientRequests() {
	// check if any users are
	client.GetAllUserFromServer()

	//create new users
	resp, err := client.AddUserToServer("Dima", 22)
	if err != nil {
		log.Printf("[CLIENT] error creating user: %v", err)
	}
	user1, err := client.GetUserFromResponse(resp)
	if err != nil {
		log.Printf("[CLIENT] error getting user id: %v", err)
	}

	resp, err = client.AddUserToServer("Kristina", 19)
	if err != nil {
		log.Printf("[CLIENT] error creating user: %v", err)
	}
	user2, err := client.GetUserFromResponse(resp)
	if err != nil {
		log.Printf("[CLIENT] error getting user id: %v", err)
	}

	resp, err = client.AddUserToServer("Vova", 28)
	if err != nil {
		log.Printf("[CLIENT] error creating user: %v", err)
	}
	user3, err := client.GetUserFromResponse(resp)
	if err != nil {
		log.Printf("[CLIENT] error getting user id: %v", err)
	}

	resp, err = client.AddUserToServer("Danila", 18)
	if err != nil {
		log.Printf("[CLIENT] error creating user: %v", err)
	}
	user4, err := client.GetUserFromResponse(resp)
	if err != nil {
		log.Printf("[CLIENT] error getting user id: %v", err)
	}

	// read all users from server
	client.GetAllUserFromServer()

	// read users one by one from server
	client.GetUserFromServer(user1.ID)
	client.GetUserFromServer(user2.ID)
	client.GetUserFromServer(user3.ID)
	client.GetUserFromServer(user4.ID)

	// try to read non-existing user
	client.GetUserFromServer(100500)

	// delete user from server and read all users from server
	client.DeleteUserFromServer(user4.ID)
	client.DeleteUserFromServer(100500)
	client.GetAllUserFromServer()

	// update user on server and read all users from server
	client.GetUserFromServer(user1.ID)
	user1.Age = 32
	client.UpdateUserAtServer(user1.ID, user1)
	client.GetUserFromServer(user1.ID)

	// create friendship and check updated user
	client.MakeFriendsOnServer(user1, user2)

	client.GetUserFromServer(user1.ID)
	client.GetFriendsFromServer(user1.ID)

	// delete user and check his friends
	client.DeleteUserFromServer(user2.ID)
	client.GetUserFromServer(user1.ID)
	client.GetAllUserFromServer()
}
