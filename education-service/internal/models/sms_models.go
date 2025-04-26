package models

type SmsTemplate struct {
	ID              int    `json:"id"`
	Texts           string `json:"texts"`
	ActionType      string `json:"action_type"`
	SmsTemplateType string `json:"sms_template_type"`
	IsActive        bool   `json:"is_active"`
}
