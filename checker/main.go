package checker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"webtester/sms"
)

type urlContext interface {
	WebChecker([]sms.SmsContact)
	checkCtxChanged() bool
}

type UrlCheck struct {
	Url           string
	Substring     string
	PeriodSeconds int
	state         bool
}

func NewUrlCheck(url, substring string, period int) UrlCheck {
	return UrlCheck{
		Url:           url,
		Substring:     substring,
		PeriodSeconds: period,
		state:         true,
	}
}

func (u *UrlCheck) checkCtxChanged(htmlBody string) bool {
	subStringFound := strings.Contains(htmlBody, u.Substring)

	if subStringFound != u.state {
		u.state = subStringFound
		return true
	}
	return false
}

func (u UrlCheck) WebChecker(contacts []sms.SmsContact) {
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

		if u.checkCtxChanged(string(bodyBytes)) {
			fmt.Printf("Item %s state changed to %t.\n", u.Url, u.state)
			if u.state == true {
				for _, contact := range contacts {
					message := fmt.Sprintf("Item %s is IN STOCK!", u.Url)
					if err := contact.Send(message); err != nil {
						log.Fatal(err)
					}
				}
			}
		} else {
			fmt.Printf("Item %s state NOT changed.\n", u.Url)
		}
	}
}
