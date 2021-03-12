package routes

import (
	"http/server/controllers"
	"net/http"
)

func SetHomeRoute() {
	homeController := controllers.HomeController{}
	http.HandleFunc("/", homeController.Home)
}
