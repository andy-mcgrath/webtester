package httpchecker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"webtester/config"
	"webtester/sms"
)

func WebChecker(u config.UrlCheck, contacts []sms.SmsContact) {

	method := "GET"

	client := &http.Client{Timeout: 5000 * time.Millisecond}
	req, err := http.NewRequest(method, u.Url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		if isItemAvailable(string(bodyBytes), u.Substring) {
			fmt.Printf("Item %s IN STOCK!\n", u.Url)
			for _, contact := range contacts {
				message := fmt.Sprintf("Item %s is IN STOCK!", u.Url)
				if err := contact.Send(message); err != nil {
					log.Fatal(err)
				}
			}
			os.Exit(0)
		} else {
			fmt.Printf("Item %s NOT in stock.\n", u.Url)
		}
	}
}

func isItemAvailable(htmlBody, subString string) bool {
	return strings.Contains(htmlBody, subString)
}
