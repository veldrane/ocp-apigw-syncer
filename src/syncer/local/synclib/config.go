package synclib

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

func GetConfig() Config {

	var config Config

	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, _ := os.ReadFile(filename)

	err := yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	config.check()
	config.setDefault()

	return config
}

func (c *Config) check() {

	if (c.Namespace == "") ||
		(c.Deployment == "") ||
		(c.HostHeader == "") ||
		(c.HttpPath == "") ||
		(c.HostDomain == "") {
		panic("No mandatory fields have been found - please check config.yaml")
	}

}

func (c *Config) setDefault() {

	if c.HttpsPort == "" {
		c.HttpsPort = "8443"
	}

	if c.Retries == 0 {
		c.Retries = 5
	}

	if c.SyncTimeout == 0 {
		c.SyncTimeout = 150
	}

	if c.ConnTimeout == 0 {
		c.ConnTimeout = 200
	}

	if c.Deadline == 0 {
		c.Deadline = 1000
	}

	if c.MaxKeepAlives == 0 {
		c.MaxKeepAlives = 512
	}

	if c.HostKeepAlives == 0 {
		c.HostKeepAlives = 64
	}

}
