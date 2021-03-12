package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"http/server/user"
)

type UserController struct{}

func (x *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := user.GetUsers()
	json.NewEncoder(w).Encode(users)
}

func (x *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	u, ok := user.GetUser(id)
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (x *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser user.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = user.CreateUser(newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (x *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedUser user.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedUser.ID = id

	err = user.UpdateUser(updatedUser)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(updatedUser)
}

func (x *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = user.DeleteUser(id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusBadRequest)
		return
	}

	http.Error(w, "User not found", http.StatusNotFound)
}
