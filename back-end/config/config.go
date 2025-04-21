package config

import (
	"flag"
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	LogFilepath      string `default:"back-end/server.log"`
	Verbose          bool   `default:"false"`
	DBFilepath       string `default:"back-end/db/db.db"`
	DBMigrationsPath string `default:"back-end/db/migrations/"`
	TemplatesPath    string `default:"front-end/"`
}

var path = flag.String("config", "", "path to the config file")

func parseConfig() Config {
	if *path == "" {
		c := Config{}
		defaults.SetDefaults(&c)
		return c
	}

	var c Config
	if _, err := toml.DecodeFile(*path, &c); err != nil {
		panic(fmt.Errorf("cannot read config file: %w", err))
	}

	defaults.SetDefaults(&c)

	return c
}

var once = sync.OnceValue(parseConfig)

func Get() Config {
	return once()
}
