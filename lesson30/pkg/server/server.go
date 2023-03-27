package server

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"httpServer/pkg/models"
	"httpServer/pkg/storage"
	"io"
	"log"
	"net/http"
	"strconv"
)

var store *storage.MemStorage

type Server struct {
	Router *chi.Mux
	// Db, config can be added here
}

type UserResponse struct {
	Status   int         `json:"status"`
	Response models.User `json:"response,omitempty"`
	Error    error       `json:"error,omitempty"`
}

type UsersResponse struct {
	Status   int           `json:"status"`
	Response []models.User `json:"response,omitempty"`
	Error    error         `json:"error,omitempty"`
}

func init() {
	store = storage.NewStorage()
}

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	r := s.Router
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)

	r.Route("/", func(r chi.Router) {
		r.Get("/", DefaultRoute)

		// sub-route 'users'
		r.Route("/users", func(r chi.Router) {
			r.With(paginate).Get("/", GetAllUsers)
			r.Post("/create", CreateUser)
			r.Post("/make_friend", MakeFriend)

			// sub-route by user's id
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", GetUser)
				r.Put("/", UpdateUser)
				r.Delete("/", DeleteUser)
			})
		})

		// sub-route 'friends'
		r.Route("/friends", func(r chi.Router) {
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", GetFriends)
			})
		})
	})
}

func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("hello, there!")); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, models.User{}, err)
		return
	}

	var u *models.User
	if err = json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, models.User{}, err)
		return
	}

	if err = store.Save(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, *u, err)
		log.Printf("[SERVER] error saving user to the storage: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	sendResponse(w, http.StatusCreated, *u, nil)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(userIDString)
	user := store.Get(userID)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		error404 := errors.New("user " + userIDString + " not found")
		sendResponse(w, http.StatusNotFound, models.User{}, error404)
		return
	}

	w.WriteHeader(http.StatusOK)
	sendResponse(w, http.StatusOK, *user, nil)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := store.AllUsers()

	if users == nil {
		w.WriteHeader(http.StatusNotFound)
		error404 := errors.New("users not found")
		sendResponse(w, http.StatusNotFound, models.User{}, error404)
		return
	}

	w.WriteHeader(http.StatusOK)
	sendResponseSlice(w, http.StatusOK, users, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(userIDString)
	err := store.Delete(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		sendResponse(w, http.StatusNotFound, models.User{}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	sendResponse(w, http.StatusOK, models.User{}, nil)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(userIDString)

	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, models.User{}, err)
		return
	}

	var u *models.User
	if err = json.Unmarshal(content, &u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, models.User{}, err)
		return
	}

	user, err := store.Update(userID, u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, *user, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	sendResponse(w, http.StatusInternalServerError, *user, nil)
}

func MakeFriend(w http.ResponseWriter, r *http.Request) {
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, models.User{}, err)
		return
	}

	var req models.FriendRequest
	// req := models.FriendRequest{1, 2}
	err = json.Unmarshal(content, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendResponse(w, http.StatusInternalServerError, models.User{}, err)
		return
	}

	source := store.Get(req.Source_id)
	target := store.Get(req.Target_id)
	if &source != nil && &target != nil {
		newSource, err := source.MakeFrienship(target)
		if err != nil {
			log.Println(err)
			return
		}
		store.Update(source.ID, newSource)
		updatedSource := store.Get(source.ID)

		w.WriteHeader(http.StatusOK)
		sendResponse(w, http.StatusOK, *updatedSource, nil)
	}
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(userIDString)
	user := store.Get(userID)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		error404 := errors.New("user " + userIDString + " not found")
		sendResponse(w, http.StatusInternalServerError, models.User{}, error404)
		return
	}

	var friends []models.User
	for _, user := range user.Friends {
		friends = append(friends, *user)
	}

	w.WriteHeader(http.StatusOK)
	sendResponseSlice(w, http.StatusOK, friends, nil)
}

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func sendResponse(w http.ResponseWriter, status int, user models.User, err error) {
	response := UserResponse{
		Status:   status,
		Response: user,
		Error:    err,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
		return
	}
}

func sendResponseSlice(w http.ResponseWriter, status int, users []models.User, err error) {
	response := UsersResponse{
		Status:   status,
		Response: users,
		Error:    err,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("[SERVER] can't prepare response: %v\n", err)
		return
	}
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("[SERVER] can't send response: %v\n", err)
		return
	}
}
