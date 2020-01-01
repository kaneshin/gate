package internal

import (
	"flag"
	"io/ioutil"
	"os"
	"sync"

	"github.com/goccy/go-yaml"
)

// Config represents a configuration of commands.
var Config = config{
	Targets: sync.Map{},
}

// Load parses flag options and loads yaml file.
func Load() error {
	var configPath string
	flag.StringVar(&configPath, "config", "$HOME/.config/gate/config.yml", "")
	flag.Parse()

	b, err := ioutil.ReadFile(os.ExpandEnv(configPath))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &Config)
	if err != nil {
		return err
	}

	err = Config.apply()
	if err != nil {
		return err
	}
	return nil
}
