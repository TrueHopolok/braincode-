package main

import (
	"io"
	"log"
	"os"

	plog "github.com/TrueHopolok/plog"
)

const LOG_FILE_NAME = "back-end/server.log"
var log_file *os.File
var logger *plog.Logger

func init() {
	log_file, err := os.OpenFile(LOG_FILE_NAME, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log_mw := io.MultiWriter(log_file, os.Stdout)
	logger, err = plog.NewLogger(plog.LevelInfo, log_mw, plog.RequireTimestamp|plog.RequireLevel, false)
	if err != nil {
		log.Fatalln(err)
	}
}