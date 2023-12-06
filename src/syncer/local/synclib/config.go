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

	return config
}
