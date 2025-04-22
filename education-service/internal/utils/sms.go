package utils

import (
	"fmt"
	"github.com/realtemirov/eskizuz"
)

func SendSMS(phoneNumber, message string) error {
	eskizClient, err := eskizuz.GetToken(&eskizuz.Auth{
		Email:    "abdulla.ergashev.2020@mail.ru",
		Password: "V16Q4KCD008jNLmgQ2zcnxxT5tgNM085BJShe17a",
	})
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	sms := &eskizuz.SMS{
		MobilePhone: phoneNumber,
		Message:     message,
		From:        "4546",
	}

	result, err := eskizClient.Send(sms)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	fmt.Printf("SMS sent successfully: %+v\n", result)
	return nil
}
