package main

import (
	"flag"
	"net/http"

	db "github.com/TrueHopolok/braincode-/server/db"
	logger "github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/prepared"
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
	defer db.Conn.Close()
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
	EnableFileHandlers()
	EnableControllerHandlers()
	go http.ListenAndServe(":8080", LoggerMiddleware(http.DefaultServeMux))
	logger.Log.Info("HTTP server: start succeeded")

	//* Console init
	ConsoleHandler()
}
