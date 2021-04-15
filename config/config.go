package config

import (
	"fmt"
	"log"
	"webtester/checker"
	"webtester/sms"

	"github.com/spf13/viper"
)

type confInterface map[string]interface{}

func LoadConfigFile(location string) (*viper.Viper, error) {
	mainViper := viper.New()
	mainViper.AddConfigPath(location)
	mainViper.SetConfigName("config")
	mainViper.SetConfigType("yaml")
	if err := mainViper.ReadInConfig(); err != nil {
		return nil, err
	}

	return mainViper, nil
}

func GetConfig(config *viper.Viper) (int, []*checker.UrlCheck) {
	contacts := make([]*sms.SmsContact, 0)
	urls := make([]*checker.UrlCheck, 0)

	smsJwt := config.GetString("smsWorks")
	checkInterval := config.GetInt("interval")
	contactConfigs := make([]confInterface, 0)
	webchecks := make([]confInterface, 0)

	if err := config.UnmarshalKey("contacts", &contactConfigs); err != nil {
		log.Fatal(err)
	}
	if err := config.UnmarshalKey("webcheck", &webchecks); err != nil {
		log.Fatal(err)
	}

	for _, v := range contactConfigs {
		name := fmt.Sprintf("%v", v["name"])
		dest := fmt.Sprintf("%v", v["destination"])

		fmt.Printf("Adding contact %s\n", name)
		contacts = append(contacts, sms.NewSmsContact(name, dest))
	}

	contact := sms.NewContacts(smsJwt, contacts)

	for _, v := range webchecks {
		url := fmt.Sprintf("%v", v["url"])
		substring := fmt.Sprintf("%v", v["substring"])

		fmt.Printf("Adding URL to monitor %s, looking for \"%s\"\n", url, substring)
		urls = append(urls, checker.NewUrlCheck(url, substring, contact))
	}

	return checkInterval, urls
}
