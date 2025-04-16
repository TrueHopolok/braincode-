// Implements initialization, global calls and destruction for cutomly written plog logger
package logger

import (
	"flag"
	"io"
	"os"

	plog "github.com/TrueHopolok/plog"
)

// Contain path to the log file
var LOG_FILE_NAME = flag.String("log_file", "back-end/server.log", "File path where logs will be saved into")

// Contain if logger output level is debug if true
var LOG_LEVEL_DEBUG = flag.Bool("log_level", true, "If true logger will output level debug")

// Logger variable that must initialized via Start() function in the package
var Log *plog.Logger

/*
Initialize global logger by opening log file with write, append and create flags.
Should be called once. Otherwise may result in unexpected behaviour.

	if debug: level = debug, output += os.Stdout
	else: level = info
	return os.OpenFile(log_file)
*/
func Init(debug bool) error {
	log_file, err := os.OpenFile(*LOG_FILE_NAME, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
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
