package main

import (
	"cmp"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/config"
	controllers "github.com/TrueHopolok/braincode-/server/controllers"
	"github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/session"
)

func MuxHTTP() http.Handler {
	mux := http.NewServeMux()
	EnableFileHandlers(mux)
	EnableControllerHandlers(mux)
	return LoggerMiddleware(mux)
}

func EnableFileHandlers(mux *http.ServeMux) {
	h := http.FileServer(http.Dir(config.Get().StaticPath))
	mux.Handle("GET /static/", http.StripPrefix("/static/", h))
	mux.Handle("GET /task/static/", http.StripPrefix("/task/static/", h))
	mux.Handle("GET /login/static/", http.StripPrefix("/login/static/", h))
	mux.Handle("GET /register/static/", http.StripPrefix("/register/static/", h))
	mux.Handle("GET /stats/static/", http.StripPrefix("/stats/static/", h))
	mux.Handle("GET /upload/static/", http.StripPrefix("/upload/static/", h))

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/static.favicon")
	})
}

func EnableControllerHandlers(mux *http.ServeMux) {
	mux.Handle("GET /api/tasks/", session.MiddlewareFunc(controllers.ProblemsAPI))
	mux.Handle("GET /api/submissions/", session.AuthMiddlewareFunc(controllers.SubmissionsAPI))

	mux.Handle("GET /", session.MiddlewareFunc(controllers.ProblemsPage))
	mux.Handle("DELETE /", session.AuthMiddlewareFunc(controllers.TaskDelete))

	mux.Handle("GET /task/", session.MiddlewareFunc(controllers.TaskPage))
	mux.Handle("POST /task/", session.AuthMiddlewareFunc(controllers.TaskSolve))

	mux.Handle("GET /login/", session.NoAuthMiddlewareFunc(controllers.LoginPage))
	mux.Handle("POST /login/", session.NoAuthMiddlewareFunc(controllers.UserLogin))
	mux.Handle("POST /logout/", session.AuthMiddlewareFunc(controllers.UserLogout))

	mux.Handle("GET /register/", session.NoAuthMiddlewareFunc(controllers.RegistrationPage))
	mux.Handle("POST /register/", session.NoAuthMiddlewareFunc(controllers.UserRegister))

	mux.Handle("GET /stats/", session.AuthMiddlewareFunc(controllers.ProfilePage))
	mux.Handle("POST /stats/delete-user/", session.AuthMiddlewareFunc(controllers.UserDelete))
	mux.Handle("POST /stats/change-password/", session.AuthMiddlewareFunc(controllers.UserChangePassword))

	mux.Handle("GET /upload/", session.AuthMiddlewareFunc(controllers.UploadPage))
	mux.Handle("POST /upload/", session.AuthMiddlewareFunc(controllers.TaskCreate))
}

func LoggerMiddleware(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("Caught panic: %v (%#v of type %T)", r, r, r)
				http.Error(w, "Go panicked :(", http.StatusInternalServerError)
				panic(r)
			}
		}()

		remote := cmp.Or(r.Header.Get("X-Forwarded-For"), "("+r.RemoteAddr+")")
		logger.Log.Debug("addr=%s met=%s url=%s | arrived", remote, r.Method, r.URL.Path)
		defer logger.Log.Debug("addr=%s met=%s url=%s | served", remote, r.Method, r.URL.Path)

		mux.ServeHTTP(w, r)
	})
}
