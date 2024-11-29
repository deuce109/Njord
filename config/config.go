package config

import (
	"github.com/BurntSushi/toml"
)

var WorldName string

var SaveName string

type paths struct {
	WorldPath string `toml:"world-name"`
	SaveName  string `toml:"save-name"`
}

func ReadConfig(configPath string) error {
	var p paths
	if _, err := toml.DecodeFile(configPath, &p); err != nil {
		return err
	}

	WorldName = p.WorldPath
	SaveName = p.SaveName

	return nil
}
