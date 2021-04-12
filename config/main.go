package config

import (
	"fmt"
	"log"
	"webtester/sms"

	"github.com/spf13/viper"
)

type UrlCheck struct {
	Url           string
	Substring     string
	PeriodSeconds int
}

type domain map[string]interface{}

func GetConfig() ([]sms.SmsContact, []UrlCheck) {
	contacts := make([]sms.SmsContact, 0)
	urls := make([]UrlCheck, 0)
	mainViper := viper.New()
	mainViper.AddConfigPath(".")
	mainViper.SetConfigName("config")
	mainViper.SetConfigType("yaml")
	if err := mainViper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	smsJwt := mainViper.GetString("smsWorks")
	contactConfigs := make([]domain, 0)
	webchecks := make([]domain, 0)

	if err := mainViper.UnmarshalKey("contacts", &contactConfigs); err != nil {
		log.Fatal(err)
	}
	if err := mainViper.UnmarshalKey("webcheck", &webchecks); err != nil {
		log.Fatal(err)
	}

	for _, v := range contactConfigs {
		name := fmt.Sprintf("%v", v["name"])
		dest := fmt.Sprintf("%v", v["destination"])
		fmt.Printf("Adding contact %s\n", name)
		contacts = append(contacts, sms.NewSmsContact(name, dest, smsJwt))
	}
	for _, v := range webchecks {
		url := fmt.Sprintf("%v", v["url"])
		substring := fmt.Sprintf("%v", v["substring"])
		fmt.Printf("Adding URL to monitor %s, looking for \"%s\"\n", url, substring)
		temp := UrlCheck{
			Url:           url,
			Substring:     substring,
			PeriodSeconds: v["periodSeconds"].(int),
		}
		urls = append(urls, temp)
	}

	return contacts, urls
}
