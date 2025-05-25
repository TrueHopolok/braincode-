package main

import (
	"net/http"

	controllers "github.com/TrueHopolok/braincode-/server/controllers"
)

func EnableFileHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/task/static/", http.StripPrefix("/task/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/login/static/", http.StripPrefix("/login/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/register/static/", http.StripPrefix("/register/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/stats/static/", http.StripPrefix("/stats/static/", http.FileServer(http.Dir("./frontend/static"))))
}

func EnableControllerHandlers() {
	http.HandleFunc("/", controllers.ProblemsetPage)
	http.HandleFunc("/task/", controllers.TaskPage)
	http.HandleFunc("/login/", controllers.LoginPage)
	http.HandleFunc("/register/", controllers.RegistrationPage)
	http.HandleFunc("/stats/", controllers.StatsPage)
}
