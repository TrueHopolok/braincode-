package main

import (
	"net/http"

	controllers "github.com/TrueHopolok/braincode-/server/controllers"
)

func EnableFileHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/task/static/", http.StripPrefix("/task/static/", http.FileServer(http.Dir("./frontend/static"))))
}

func EnableControllerHandlers() {
	http.HandleFunc("/", controllers.Problemset)
	http.HandleFunc("/task/", controllers.Taskpage)
}
