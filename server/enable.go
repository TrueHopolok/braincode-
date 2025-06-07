package main

import (
	"net/http"

	controllers "github.com/TrueHopolok/braincode-/server/controllers"
	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/session"
)

func EnableFileHandlers() {
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("GET /task/static/", http.StripPrefix("/task/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("GET /login/static/", http.StripPrefix("/login/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("GET /register/static/", http.StripPrefix("/register/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("GET /stats/static/", http.StripPrefix("/stats/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.Handle("GET /upload/static/", http.StripPrefix("/upload/static/", http.FileServer(http.Dir("./frontend/static"))))
	http.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/static.favicon")
	})
}

func EnableControllerHandlers() {
	http.Handle("GET /", session.MiddlewareFunc(controllers.ProblemsPage))
	http.Handle("DELETE /", session.AuthMiddlewareFunc(controllers.TaskDelete))

	http.Handle("GET /task/", session.MiddlewareFunc(controllers.TaskPage))
	http.Handle("POST /task/", session.AuthMiddlewareFunc(controllers.TaskSolve))

	http.Handle("GET /login/", session.NoAuthMiddlewareFunc(controllers.LoginPage))
	http.Handle("POST /login/", session.NoAuthMiddlewareFunc(controllers.UserLogin))
	http.Handle("DELETE /login/", session.AuthMiddlewareFunc(controllers.UserLogout))

	http.Handle("GET /register/", session.NoAuthMiddlewareFunc(controllers.RegistrationPage))
	http.Handle("POST /register/", session.NoAuthMiddlewareFunc(controllers.UserRegister))

	http.Handle("GET /stats/", session.AuthMiddlewareFunc(controllers.StatsPage))
	http.Handle("DELETE /stats/", session.AuthMiddlewareFunc(controllers.UserDelete))

	http.Handle("GET /upload/", session.AuthMiddlewareFunc(controllers.UploadPage))
	http.Handle("POST /upload/", session.AuthMiddlewareFunc(controllers.TaskCreate))
}

func LoggerMiddleware(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("req=%p met=%s url=%s | arrived", r, r.Method, r.URL.Path)
		defer logger.Log.Debug("req=%p met=%s url=%s | served", r, r.Method, r.URL.Path)

		mux.ServeHTTP(w, r)
	})
}
