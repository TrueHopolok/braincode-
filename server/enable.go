package main

import (
	"net/http"

	controllers "github.com/TrueHopolok/braincode-/server/controllers"
	"github.com/TrueHopolok/braincode-/server/session"
)

func EnableFileHandlers() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/task/static/", http.StripPrefix("/task/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/login/static/", http.StripPrefix("/login/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/register/static/", http.StripPrefix("/register/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/stats/static/", http.StripPrefix("/stats/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("/upload/static/", http.StripPrefix("/upload/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/static.favicon")
	})
}

func EnableControllerHandlers() {
	http.Handle("GET /", session.MiddlewareFunc(controllers.ProblemsetPage))
	http.Handle("DELETE /", session.AuthMiddlewareFunc(controllers.ProblemsetPage))

	http.Handle("GET /task/", session.MiddlewareFunc(controllers.TaskPage))
	http.Handle("POST /task/", session.AuthMiddlewareFunc(controllers.TaskPage))

	http.Handle("/login/", session.NoAuthMiddlewareFunc(controllers.LoginPage))
	http.Handle("DELETE /login/", session.AuthMiddlewareFunc(controllers.LoginPage))
	http.Handle("/register/", session.NoAuthMiddlewareFunc(controllers.RegistrationPage))

	http.Handle("/stats/", session.AuthMiddlewareFunc(controllers.StatsPage))
	http.Handle("/upload/", session.AuthMiddlewareFunc(controllers.UploadPage))
}
