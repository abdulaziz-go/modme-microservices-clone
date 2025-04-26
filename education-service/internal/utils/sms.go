package utils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/lib/pq"
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
		fmt.Println("SMS failed:", string(body))
		//return errors.New(string(body))
	}

	fmt.Println("SMS sent:", string(body))
	return nil
}

func GetSmsFormatted(sms, teacher string, db *sql.DB, studentID, groupID string, amountValue float64, companyId string) (string, int) {
	var studentName string
	err := db.QueryRow("SELECT name FROM students WHERE id = $1", studentID).Scan(&studentName)
	if err != nil {
		fmt.Println("Error fetching student:", err)
		studentName = "(Student)"
	}

	var groupName, startTime string
	var days []string
	var roomID, companyID int
	err = db.QueryRow(`SELECT name, start_time, days, room_id, company_id FROM groups WHERE id = $1`, groupID).Scan(&groupName, &startTime, pq.Array(&days), &roomID, &companyID)
	if err != nil {
		fmt.Println("Error fetching group:", err)
		groupName = "(Group)"
		startTime = "(Time)"
	}

	var roomName string
	err = db.QueryRow("SELECT title FROM rooms WHERE id = $1", roomID).Scan(&roomName)
	if err != nil {
		fmt.Println("Error fetching room:", err)
		roomName = "(Room)"
	}

	var companyName string
	err = db.QueryRow("SELECT title FROM company WHERE id = $1", companyId).Scan(&companyName)
	if err != nil {
		fmt.Println("Error fetching company:", err)
		companyName = "(LC)"
	}

	fmt.Println("here company name from response ", companyName)

	daysStr := strings.Join(days, ", ")

	replacements := map[string]string{
		"(STUDENT)": studentName,
		"(GROUP)":   groupName,
		"(TIME)":    startTime,
		"(LC)":      companyName,
		"(TEACHER)": teacher,
		"(DAYS)":    daysStr,
		"(ROOM)":    roomName,
		"(SUM)":     fmt.Sprintf("%.2f", amountValue),
	}

	for key, value := range replacements {
		sms = strings.ReplaceAll(sms, key, value)
	}
	length := len([]rune(sms))

	isCyrillic := false
	for _, r := range sms {
		if r >= 0x0400 && r <= 0x04FF {
			isCyrillic = true
			break
		}
	}

	var smsCount int
	if isCyrillic {
		smsCount = (length + 69) / 70
	} else {
		smsCount = (length + 159) / 160
	}

	return sms, smsCount
}
