package controllers

import (
	"fmt"
	"net/http"
)

type HomeController struct{}

func (x *HomeController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
