package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

var smsToken string

func getEskizToken() (string, error) {
	url := "https://notify.eskiz.uz/api/auth/login"
	payload := strings.NewReader("email=abdulla.ergashev.2020@mail.ru&password=V16Q4KCD008jNLmgQ2zcnxxT5tgNM085BJShe17a")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	token := data["token"].(string)

	return token, nil
}

func SendSMS(phoneNumber, message string) error {
	if smsToken == "" {
		token, err := getEskizToken()
		if err != nil {
			return err
		}
		smsToken = token
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	w.WriteField("from", "4546")
	w.WriteField("mobile_phone", strings.TrimPrefix(phoneNumber, "+"))
	w.WriteField("message", message)
	w.Close()

	req, err := http.NewRequest("POST", "https://notify.eskiz.uz/api/message/sms/send", &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+smsToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		// retry once
		token, err := getEskizToken()
		if err != nil {
			return err
		}
		smsToken = token
		return SendSMS(phoneNumber, message)
	}

	fmt.Println("SMS sent:", string(body))
	return nil
}
