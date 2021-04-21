package user

import (
	"fmt"
	"sync"
)

var (
	globalUserIDCountMtx = sync.Mutex{}
	globalUserIDCount    = int(0)
	globalUsers          = []User{}
)

func getNextID() int {
	globalUserIDCountMtx.Lock()
	defer globalUserIDCountMtx.Unlock()
	globalUserIDCount++
	return globalUserIDCount
}

func findUserByID(id int) int {
	for i, v := range globalUsers {
		if id == v.ID {
			return i
		}
	}
	return -1
}

func GetUsers() []User {
	users := []User{}
	copy(users, globalUsers)
	return users
}

func GetUser(id int) (User, bool) {
	idx := findUserByID(id)
	if idx == -1 {
		return User{}, false
	}
	return globalUsers[idx], true
}

func CreateUser(user User) error {
	idx := findUserByID(user.ID)
	if idx != -1 {
		return fmt.Errorf("user id exists")
	}
	user.ID = getNextID()
	globalUsers = append(globalUsers, user)
	return nil
}

func UpdateUser(user User) error {
	idx := findUserByID(user.ID)
	if idx == -1 {
		return fmt.Errorf("user id not exists")
	}
	globalUsers[idx] = user
	return nil
}

func DeleteUser(id int) error {
	idx := findUserByID(id)
	if idx == -1 {
		return fmt.Errorf("user id not exists")
	}
	globalUsers = append(globalUsers[:idx], globalUsers[idx+1:]...)
	return nil
}
