package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	logger "github.com/TrueHopolok/braincode-/back-end/logger"
)

func main() {
	err := logger.Start(true)
	if err != nil {
		log.Fatalln(err)
	}
	logger.Log.Line()
	defer logger.Stop()
	
	http.HandleFunc("/", HelloServer)
	go http.ListenAndServe(":8080", nil)
	
	logger.Log.Info("Server started")

	ConsoleHandler()
}

func StopServer() {
	logger.Log.Info("Server stopped, via console input")
	logger.Stop()
	os.Exit(0)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("Request=%p arrived", r)
	defer logger.Log.Debug("Request=%p served", r)
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}