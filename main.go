package main

import (
	"time"
	"webtester/config"
)

func main() {

	contacts, urls := config.GetConfig()

	for {
		for _, url := range urls {
			go url.WebChecker(contacts)
		}
		time.Sleep(300 * time.Second)
	}
}
