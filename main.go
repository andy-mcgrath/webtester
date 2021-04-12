package main

import (
	"time"
	"webtester/config"
)

func main() {

	interval, contacts, urls := config.GetConfig()

	for {
		for _, url := range urls {
			go url.WebChecker(contacts)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
