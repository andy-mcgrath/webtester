package main

import (
	"time"
	"webtester/config"
	"webtester/httpchecker"
)

func main() {

	contacts, urls := config.GetConfig()

	for {
		for _, url := range urls {
			go httpchecker.WebChecker(url, contacts)
		}
		time.Sleep(300 * time.Second)
	}
}
