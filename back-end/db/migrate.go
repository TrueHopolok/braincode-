package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/TrueHopolok/braincode-/back-end/config"
	"github.com/TrueHopolok/braincode-/back-end/logger"
)

/*
Execute all migrations with numbered prefixes.
Starting from "001_" prefix and adding 1 to the next prefix.

	max_prefix = "999_"
	if cur_prefix do not exist: migration execution stops
*/
func Migrate() error {
	files, err := os.ReadDir(config.Get().DBMigrationsPath)
	if err != nil {
		return err
	}
	for i := 1; i <= 999; i++ {
		prefix := fmt.Sprintf("%03d_", i)
		found := false
		for _, file := range files {
			fname := file.Name()
			if !file.IsDir() && strings.HasPrefix(fname, prefix) {
				logger.Log.Debug("Migration: found = %s", fname)
				buf, err := os.ReadFile(fname)
				_, err = Conn.Exec(string(buf))
				if err != nil {
					return err
				}
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return nil
}
