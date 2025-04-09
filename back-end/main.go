package main

import (
	"fmt"
	"os"

	plog "github.com/TrueHopolok/plog"
)

const LOG_FILE_NAME = "back-end/server.log"

func main() {
	flog, err := os.OpenFile(LOG_FILE_NAME, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer flog.Close()
	logger, _ := plog.NewLogger(plog.LevelInfo, flog, plog.RequireTimestamp|plog.RequireLevel, false)
	logger.Info("program started / logger initializated")
}
