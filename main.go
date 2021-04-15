package main

import (
	"flag"
	"log"
	"time"
	"webtester/config"
)

func main() {
	var configLocation string
	flag.StringVar(&configLocation, "config", ".", "Config file 'config.yaml' location")
	flag.Parse()

	configFile, err := config.LoadConfigFile(configLocation)
	if err != nil {
		log.Fatal(err)
	}

	interval, urls := config.GetConfig(configFile)

	for {
		for _, url := range urls {
			go url.WebChecker()
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
