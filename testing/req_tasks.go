package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		log.Fatalln("1", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("2", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("3", err)
	}
	fmt.Println(string(b))
}
