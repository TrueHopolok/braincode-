package main

import (
	"flag"
	"net/http"

	"github.com/TrueHopolok/braincode-/server/config"
	db "github.com/TrueHopolok/braincode-/server/db"
	logger "github.com/TrueHopolok/braincode-/server/logger"
	"github.com/TrueHopolok/braincode-/server/prepared"
)

func main() {
	flag.Parse()

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

	// Error channel for all concurent threads
	httpChan := make(chan error)

	//* HTTP init
	logger.Log.Info("HTTP server: starting...")
	go func() {
		httpChan <- http.ListenAndServe(":8080", MuxHTTP())
	}()
	select {
	case err := <-httpChan:
		logger.Log.Fatal("HTTP server: start failed; error=%s", err)
	default:
		logger.Log.Info("HTTP server: start succeeded")
	}

	//* Console init
	consoleChan := make(chan error)
	quitChan := make(chan bool)
	if config.Get().EnableConsole {
		go func() {
			consoleChan <- ConsoleHandler(quitChan)
		}()
	}

	select {
	case err := <-httpChan:
		logger.Log.Fatal("HTTP server: execution stopped; err=%s", err)
	case err := <-consoleChan:
		logger.Log.Fatal("Console: stopped execution; err=%s", err)
	case <-quitChan:
		logger.Log.Warn("Console: quitChan returned a value, server closing")
	}
}
