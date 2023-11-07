package services

import "os"

type SMSService interface {
	SetAPIKey(apikey string)
	SendSMS(number string, message string) error
}

// var SMS *KavehnegarSMSService
var SMS *GhasedakSMSService

func RegisterSMSService() {
	//var sms KavehnegarSMSService
	//sms.SetAPIKey(os.Getenv("KAVEHNEGAR_API_KEY")
	//SMS = &sms

	var sms GhasedakSMSService
	sms.SetAPIKey(os.Getenv("GHASEDAK_API_KEY"))
	SMS = &sms
}
