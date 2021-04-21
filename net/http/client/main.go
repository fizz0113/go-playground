package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "http://localhost:8080"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func printResponse(action string, body io.Reader) {
	fmt.Println("===", action, "===")
	data, _ := io.ReadAll(body)
	fmt.Println(string(data))
	fmt.Println()
}

func getUsers() {
	resp, err := http.Get(baseURL + "/users")
	if err != nil {
		fmt.Println("GET /users error:", err)
		return
	}
	defer resp.Body.Close()
	printResponse("GET /users", resp.Body)
}

func getUser(id int) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, id))
	if err != nil {
		fmt.Println("GET /users/{id} error:", err)
		return
	}
	defer resp.Body.Close()
	printResponse(fmt.Sprintf("GET /users/%d", id), resp.Body)
}

func createUser(name string) {
	user := User{Name: name}
	body, _ := json.Marshal(user)
	resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("POST /users error:", err)
		return
	}
	defer resp.Body.Close()
	printResponse("POST /users", resp.Body)
}

func updateUser(id int, newName string) {
	updated := User{Name: newName}
	body, _ := json.Marshal(updated)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/users/%d", baseURL, id), bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("PUT /users/{id} error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("PUT request failed:", err)
		return
	}
	defer resp.Body.Close()
	printResponse(fmt.Sprintf("PUT /users/%d", id), resp.Body)
}

func deleteUser(id int) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d", baseURL, id), nil)
	if err != nil {
		fmt.Println("DELETE /users/{id} error:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("DELETE request failed:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("DELETE /users/%d: Status Code = %d\n", id, resp.StatusCode)
}

func main() {
	// 測試：取得所有使用者
	getUsers()

	// 測試：新增使用者
	createUser("Charlie")

	// 測試：取得 ID 為 1 的使用者
	getUser(1)

	// 測試：更新 ID 為 1 的使用者名稱
	updateUser(1, "Updated Alice")

	// 測試：刪除 ID 為 2 的使用者
	deleteUser(2)
}
