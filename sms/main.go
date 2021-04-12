package sms

import (
	"fmt"
	"net/http"
	"strings"
)

type SmsSenter interface {
	Send() error
}

type SmsContact struct {
	name        string
	destination string
	*smsAuth
}

func NewSmsContact(name, destination, smsJwt string) SmsContact {
	return SmsContact{
		name:        name,
		destination: destination,
		smsAuth:     SetInstance(smsJwt),
	}
}

func (sms *SmsContact) Send(message string) error {

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
	req.Header.Add("Authorization", sms.smsAuth.jwt)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fmt.Printf("SMS sent to %s, status %s\n", sms.name, res.Status)

	return nil
}
