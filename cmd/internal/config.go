package internal

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents a configuration of commands.
var Config = struct {
	Gate struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"gate"`
	Slack struct {
		App struct {
			Incoming []struct {
				URL     string `toml:"url"`
				Channel string `toml:"channel"`
			}
		} `toml:"app"`
		Incoming struct {
			URL       string `toml:"url"`
			Channel   string `toml:"channel"`
			Username  string `toml:"username"`
			IconEmoji string `toml:"icon_emoji"`
		} `toml:"incoming"`
	} `toml:"slack"`
	LINE struct {
		Notify struct {
			AccessToken string `toml:"access_token"`
		} `toml:"notify"`
	} `toml:"line"`
	Facebook struct {
		Messenger struct {
			ID          string `toml:"id"`
			AccessToken string `toml:"access_token"`
		} `toml:"messenger"`
	} `toml:"facebook"`
}{}

// ParseFlag parses flag options and toml file.
func ParseFlag() error {
	var configPath string
	flag.StringVar(&configPath, "config", "$HOME/.config/gate.tml", "")
	flag.Parse()

	_, err := toml.DecodeFile(os.ExpandEnv(configPath), &Config)
	return err
}
