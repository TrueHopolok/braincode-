// Implements initialization, global calls and destruction for custom plog logger
package logger

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"io"
	"log"
	"os"

	"github.com/TrueHopolok/braincode-/server/config"
	plog "github.com/TrueHopolok/plog"
)

// Logger variable that must initialized via Start() function in the package
var Log *plog.Logger

func Start() {
	path := config.Get().LogFilepath
	log_file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log_writer := io.Writer(log_file)
	log_level := plog.LevelInfo
	if config.Get().Verbose {
		log_writer = io.MultiWriter(log_file, os.Stdout)
		log_level = plog.LevelDebug
	}
	Log, err = plog.NewLogger(log_level, log_writer, plog.RequireTimestamp|plog.RequireLevel, false)
	if err != nil {
		log.Fatalln(err)
	}
}
