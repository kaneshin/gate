package internal

import (
	"strings"
	"sync"
)

const (
	slackIncoming = "slack.incoming."
	lineNotify    = "line.notify."
)

type (
	slack struct {
		Incoming map[string]string `yaml:"incoming"`
	}

	line struct {
		Notify map[string]string `yaml:"notify"`
	}

	config struct {
		Targets       sync.Map
		DefaultTarget string `yaml:"default_target"`
		Env           struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"env"`
		Slack slack `yaml:"slack"`
		LINE  line  `yaml:"line"`
	}
)

func (c *config) apply() error {
	for k, v := range c.Slack.Incoming {
		c.Targets.Store(slackIncoming+k, v)
	}
	for k, v := range c.LINE.Notify {
		c.Targets.Store(lineNotify+k, v)
	}
	return nil
}

func (c slack) IsIncoming(name string) bool {
	return strings.HasPrefix(name, slackIncoming)
}

func (c line) IsNotify(name string) bool {
	return strings.HasPrefix(name, lineNotify)
}
