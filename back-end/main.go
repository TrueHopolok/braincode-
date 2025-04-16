package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/TrueHopolok/braincode-/back-end/db"
	logger "github.com/TrueHopolok/braincode-/back-end/logger"
)

func main() {
	var err error

	//* Logger init
	err = logger.Start(true)
	if err != nil {
		log.Fatalln(err)
	}
	logger.Log.Line()
	defer logger.Stop()

	//* Database init
	logger.Log.Info("Database: connecting...")
	if err = db.Init(); err != nil {
		logger.Log.Error("Database: connection failed; err=%s", err)
		logger.Stop()
		os.Exit(1)
	}
	logger.Log.Info("Database: connection succeeded")

	//* Database migrate
	logger.Log.Info("Migrations: executing...")
	if err = db.Migrate("drop", "create"); err != nil {
		logger.Log.Error("Migration: execution failed; error=%s", err)
		logger.Stop()
		os.Exit(1)
	}
	logger.Log.Info("Migrations: execution succeeded")

	//* HTTP init
	logger.Log.Info("HTTP server: starting...")
	http.HandleFunc("/", HelloServer) // TODO: normal function handler for complicated requests
	go http.ListenAndServe(":8080", nil)
	logger.Log.Info("HTTP server: start succeeded")

	//* Console init
	ConsoleHandler()
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("Request=%p arrived", r)
	defer logger.Log.Debug("Request=%p served", r)
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
