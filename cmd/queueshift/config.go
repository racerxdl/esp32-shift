package main

import (
	"github.com/BurntSushi/toml"
	"github.com/mewkiz/pkg/osutil"
	"github.com/racerxdl/esp32-shift/pkg/queue"
	"os"
)

type Config struct {
	queue.MQTTConfig
	SerialPort    string
	BaseTopic     string
	StartsWithOne bool
}

const configFile = "espshift.toml"

var config Config

func LoadConfig() {
	log.Info("Loading config %s", configFile)
	if !osutil.Exists(configFile) {
		log.Error("Config file %s does not exists.", configFile)
		os.Exit(1)
	}

	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Error("Error decoding file %s: %s", configFile, err)
		os.Exit(1)
	}
}
