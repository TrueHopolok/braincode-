package logger

import (
	"io"
	"os"

	plog "github.com/TrueHopolok/plog"
)

const LOG_FILE_NAME = "back-end/server.log"
var log_file *os.File
var Log *plog.Logger

func Start(debug bool) error {
	log_file, err := os.OpenFile(LOG_FILE_NAME, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	var log_writer io.Writer = log_file
	log_level := plog.LevelInfo
	if debug {
		log_writer = io.MultiWriter(log_file, os.Stdout)
		log_level = plog.LevelDebug
	}
	Log, err = plog.NewLogger(log_level, log_writer, plog.RequireTimestamp|plog.RequireLevel, false)
	return err
}

func Stop() {
	log_file.Close()
}