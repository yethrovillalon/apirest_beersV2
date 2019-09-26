package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// ReadConfig ... < Some lines that describe your function>
func ReadConfig() Config {
	var configfile = "properties.ini"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
