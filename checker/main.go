package checker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Icontact interface {
	Send(string) error
}

type UrlCheck struct {
	Url       string
	Substring string
	state     bool
	once      sync.Once
	Icontact
}

func NewUrlCheck(url, substring string, contact Icontact) *UrlCheck {
	return &UrlCheck{
		Url:       url,
		Substring: substring,
		state:     false,
		Icontact:  contact,
	}
}

func (u *UrlCheck) checkCtxChanged(htmlBody string) bool {
	subStringFound := strings.Contains(htmlBody, u.Substring)

	u.once.Do(
		func() {
			fmt.Printf("Initialising state for %s to %t\n", u.Url, subStringFound)
			u.state = subStringFound
		})

	if subStringFound != u.state {
		u.state = subStringFound
		return true
	}
	return false
}

func (u *UrlCheck) WebChecker() {
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
			if u.state {
				message := fmt.Sprintf("Item %s is IN STOCK!", u.Url)
				if err := u.Icontact.Send(message); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Printf("Item %s state NOT changed.\n", u.Url)
		}
	}
}
