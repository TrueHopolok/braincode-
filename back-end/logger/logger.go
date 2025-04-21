// Implements initialization, global calls and destruction for cutomly written plog logger
package logger

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"io"
	"log"
	"os"

	"github.com/TrueHopolok/braincode-/back-end/config"
	plog "github.com/TrueHopolok/plog"
)

// Logger variable that must initialized via Start() function in the package
var Log *plog.Logger

func init() {
	log_file, err := os.OpenFile(config.LOG_FILE_PATH, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	var log_writer io.Writer = log_file
	log_level := plog.LevelInfo
	if config.LOG_IS_DEBUG {
		log_writer = io.MultiWriter(log_file, os.Stdout)
		log_level = plog.LevelDebug
	}
	Log, err = plog.NewLogger(log_level, log_writer, plog.RequireTimestamp|plog.RequireLevel, false)
	if err != nil {
		log.Fatalln(err)
	}
}
