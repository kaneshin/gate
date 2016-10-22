package internal

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	configPath = flag.String("config", "$HOME/.gate.tml", "")

	// Config represents a configuration of commands.
	Config = struct {
		Gate struct {
			Host string `toml:"host"`
			Port int    `toml:"port"`
		} `toml:"gate"`
		Slack struct {
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
)

func ParseFlag() {
	flag.Parse()

	fp := os.ExpandEnv(*configPath)
	if _, err := toml.DecodeFile(fp, &Config); err != nil {
		log.Fatal(err)
		return
	}
}
