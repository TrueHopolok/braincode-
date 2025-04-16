package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	defer log_file.Close()
	
	http.HandleFunc("/", HelloServer)
	go http.ListenAndServe(":8080", nil)
	
	logger.Info("Server started")

	ConsoleHandler()
}

func StopServer() {
	logger.Info("Server stopped, via console input")
	os.Exit(0)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	logger.Info("Request=%p arrived", r)
	defer logger.Info("Request=%p served", r)
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}