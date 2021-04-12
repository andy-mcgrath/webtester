package sms

import (
	"fmt"
	"sync"
)

var once sync.Once

type smsAuth struct {
	jwt string
}

var singleInstance *smsAuth

func SetInstance(auth string) *smsAuth {
	if singleInstance == nil {
		once.Do(
			func() {
				fmt.Println("Creting SMS Works Singleton Instance Now")
				singleInstance = &smsAuth{
					jwt: auth,
				}
			})
	} else {
		fmt.Println("SMS Works Singleton Instance already created")
	}
	return singleInstance
}

func GetInstance() *smsAuth {
	return singleInstance
}
