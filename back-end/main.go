package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"

	db "github.com/TrueHopolok/braincode-/back-end/db"
	logger "github.com/TrueHopolok/braincode-/back-end/logger"
	"github.com/TrueHopolok/braincode-/back-end/prepared"
)

func init() {
	flag.Parse()
}

func main() {
	//* Logger init
	logger.Start()

	//* Database init
	logger.Log.Info("Database: connecting...")
	if err := db.Init(); err != nil {
		logger.Log.Fatal("Database: connection failed; err=%s", err)
	}
	logger.Log.Info("Database: connection succeeded")

	//* Database migrate
	logger.Log.Info("Migrations: executing...")
	if err := db.Migrate(); err != nil {
		logger.Log.Fatal("Migration: execution failed; error=%s", err)
	}
	logger.Log.Info("Migrations: execution succeeded")

	//* Templates init
	logger.Log.Info("Templates: initilizating...")
	if err := prepared.Init(); err != nil {
		logger.Log.Fatal("Templates: initilization failed; error=%s", err)
	}
	logger.Log.Info("Templates: initilization succeeded")

	//* HTTP init
	logger.Log.Info("HTTP server: starting...")
	http.Handle("/front-end/", http.StripPrefix("/front-end/", http.FileServer(http.Dir("./front-end"))))
	http.HandleFunc("/", indexHandler)
	go http.ListenAndServe(":8080", nil)
	logger.Log.Info("HTTP server: start succeeded")

	//* Console init
	ConsoleHandler()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("Request=%p arrived", r)
	defer logger.Log.Debug("Request=%p served", r)
	buf := bufio.NewWriter(w)
	err := prepared.Templates.ExecuteTemplate(buf, "index.html", struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	})
	if err != nil {
		fmt.Fprint(w, "ERROR")
		logger.Log.Error("Request=%p failed; error=%s", r, err)
		return
	}
	if err = buf.Flush(); err != nil {
		fmt.Fprint(w, "ERROR")
		logger.Log.Error("Request=%p failed; error=%s", r, err)
		return
	}
}
