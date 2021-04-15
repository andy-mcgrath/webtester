package config

import (
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	_, err := LoadConfigFile("..")
	if err != nil {
		t.Errorf("Loading project example config file failed with, %s", err)
	}

	viperConfig, err := LoadConfigFile("./invalid/path/")
	if err == nil {
		t.Logf("Viper Config returned, %t", viperConfig.IsSet("interval"))
		t.Errorf("Configuration loaded for invalid path, %s", err)
	}

}
