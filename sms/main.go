package sms

import (
	"fmt"
	"net/http"
	"strings"
)

type SmsSenter interface {
	Send(string) error
}

type Contacts struct {
	smsContact []*SmsContact
	*smsAuth
}

type SmsContact struct {
	name        string
	destination string
}

func NewContacts(smsJwt string, contacts []*SmsContact) *Contacts {
	return &Contacts{
		smsAuth:    SetInstance(smsJwt),
		smsContact: contacts,
	}
}

func NewSmsContact(name, destination string) *SmsContact {
	return &SmsContact{
		name:        name,
		destination: destination,
	}
}

func (c *Contacts) Send(message string) error {
	for _, contact := range c.smsContact {
		if err := contact.Send(c.smsAuth.jwt, message); err != nil {
			return err
		}
	}
	return nil
}

func (sms *SmsContact) Send(auth, message string) error {
	url := "https://api.thesmsworks.co.uk/v1/message/send"
	method := "POST"

	jsonStr := fmt.Sprintf(`{
    "sender": "ALDIcheker",
    "destination": "%s",
    "content": "Hey %s, %s"
}`, sms.destination, sms.name, message)

	payload := strings.NewReader(jsonStr)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Printf("SMS sent to %s, status %s\n", sms.name, res.Status)

	return nil
}
