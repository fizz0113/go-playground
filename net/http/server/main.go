package main

import (
	"fmt"
	"http/server/routes"
	"net/http"
)

const (
	serverHost = "localhost:8080"
)

func main() {
	routes.SetHomeRoute()
	routes.SetUserRoute()

	fmt.Println("host:", serverHost)
	http.ListenAndServe(serverHost, nil)
}
