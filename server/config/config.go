package config

import (
	"flag"
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	Verbose       bool   `default:"true"`
	EnableConsole bool   `default:"true"`
	LogFilepath   string `default:"server/server.log"`
	TemplatesPath string `default:"frontend/"`
	StaticPath    string `default:"./frontend/static"`
	DBuser        string `default:"root"`
	DBpass        string `default:"root"`
	DBname        string `default:"braincode"`
	DBqueriesPath string `default:"server/db/queries/"`
}

var CfgPath = flag.String("config", "", "path to the config file")

func parseConfig() Config {
	if *CfgPath == "" {
		c := Config{}
		defaults.SetDefaults(&c)
		return c
	}

	var c Config

	defaults.SetDefaults(&c)
	if _, err := toml.DecodeFile(*CfgPath, &c); err != nil {
		panic(fmt.Errorf("cannot read config file: %w", err))
	}

	return c
}

var once = sync.OnceValue(parseConfig)

func Get() Config {
	return once()
}
