package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	db "github.com/TrueHopolok/braincode-/back-end/db"
	logger "github.com/TrueHopolok/braincode-/back-end/logger"
)

var DEBUG_MODE bool

func init() {
	flag.BoolVar(&DEBUG_MODE, "debug", true, "Launch program in the debug mode")
}

func main() {
	var err error

	//* Logger init
	err = logger.Init(DEBUG_MODE)
	if err != nil {
		log.Fatalln(err)
	}
	logger.Log.Line()

	//* Database init
	logger.Log.Info("Database: connecting...")
	if err = db.Init(); err != nil {
		logger.Log.Fatal("Database: connection failed; err=%s", err)
	}
	logger.Log.Info("Database: connection succeeded")

	//* Database migrate
	logger.Log.Info("Migrations: executing...")
	if err = db.Migrate("drop", "create"); err != nil {
		logger.Log.Fatal("Migration: execution failed; error=%s", err)
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
