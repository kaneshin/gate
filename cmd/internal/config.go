package internal

// Config represents a configuration of commands.
var Config = config{
	Targets: map[string]string{},
}

type config struct {
	Targets       map[string]string
	DefaultTarget string `yaml:"default_target"`
	Env           struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"env"`
	Slack struct {
		Incoming map[string]string `yaml:"incoming"`
	} `yaml:"slack"`
	LINE struct {
		Notify map[string]string `yaml:"notify"`
	} `yaml:"line"`
}

func (c *config) apply() error {
	for k, v := range c.Slack.Incoming {
		c.Targets["slack.incoming."+k] = v
	}
	for k, v := range c.LINE.Notify {
		c.Targets["line.notify."+k] = v
	}
	return nil
}
