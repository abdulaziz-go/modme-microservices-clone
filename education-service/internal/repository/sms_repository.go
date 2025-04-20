package repository

import (
	"database/sql"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"fmt"
)

type SmsRepository struct {
	db *sql.DB
}

func NewSmsRepository(db *sql.DB) *SmsRepository {
	return &SmsRepository{
		db: db,
	}
}

func (r SmsRepository) GetSmsLog(companyId string, studentId string, pageReq *pb.PageRequest) (*pb.GetSmsLogResponse, error) {
	offset := (pageReq.Page - 1) * pageReq.Size

	var (
		baseQuery string
		args      []interface{}
	)

	if studentId != "" {
		baseQuery = `SELECT sms_count, array_to_string(texts, ' '), created_by_name, created_at
		             FROM sms_used WHERE company_id = $1 AND student_id = $2
		             ORDER BY created_at DESC LIMIT $3 OFFSET $4`
		args = []interface{}{companyId, studentId, pageReq.Size, offset}
	} else {
		baseQuery = `SELECT sms_count, array_to_string(texts, ' '), created_by_name, created_at
		             FROM sms_used WHERE company_id = $1
		             ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{companyId, pageReq.Size, offset}
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*pb.SmsLogList
	for rows.Next() {
		var log pb.SmsLogList
		if err := rows.Scan(&log.SmsCount, &log.SmsValue, &log.CreatorName, &log.SendDate); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	return &pb.GetSmsLogResponse{
		Datas: logs,
	}, nil
}

func (r SmsRepository) AddSms(req *pb.AddSmsRequest) (*pb.AbsResponse, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	_, err = tx.Exec(`INSERT INTO sms_payments (company_id, comment, sum, sms_count) VALUES ($1, $2, $3, $4)`,
		req.CompanyId, req.Comment, req.Sum, req.SmsCount)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`UPDATE company SET sms_balance = sms_balance + $1 WHERE id = $2`,
		req.SmsCount, req.CompanyId)
	if err != nil {
		return nil, err
	}

	return &pb.AbsResponse{
		Status:  200,
		Message: "ok",
	}, nil
}

func (r SmsRepository) DeleteSms(paymentId string) (*pb.AbsResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	var (
		smsCount  float64
		companyId int
	)

	err = tx.QueryRow(`SELECT sms_count, company_id FROM sms_payments WHERE id = $1`, paymentId).Scan(&smsCount, &companyId)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`DELETE FROM sms_payments WHERE id = $1`, paymentId)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`UPDATE company SET sms_balance = sms_balance - $1 WHERE id = $2`, smsCount, companyId)
	if err != nil {
		return nil, err
	}

	return &pb.AbsResponse{
		Status:  200,
		Message: "deleted",
	}, nil
}

func (r SmsRepository) GetSmsTransactionDetail(page int32, size int32, companyId string) (*pb.GetSmsTransactionDetailResponse, error) {
	offset := (page - 1) * size

	rows, err := r.db.Query(`
		SELECT comment, sms_count, sum , id
		FROM sms_payments
		WHERE company_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, companyId, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datas []*pb.GetSmsTransactionList
	for rows.Next() {
		var item pb.GetSmsTransactionList
		if err := rows.Scan(&item.Comment, &item.SmsCount, &item.Sum, &item.TransactionId); err != nil {
			return nil, err
		}
		datas = append(datas, &item)
	}

	return &pb.GetSmsTransactionDetailResponse{
		Datas: datas,
	}, nil
}

func (r SmsRepository) GetSmsTemplate(smsType string, companyId string) (*pb.GetSmsTemplateResponse, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if smsType == "ACTION" {
		rows, err = r.db.Query(`
			SELECT action_type, array_to_string(texts, ' '), sms_count, is_active, insufficient_balance_send_count
			FROM sms_template
			WHERE company_id = $1 AND sms_template_type = 'ACTION'
		`, companyId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var (
			datas                          []*pb.SmsTemplateList
			insufficientBalanceResendCount int32
		)

		for rows.Next() {
			var (
				item        pb.SmsTemplateList
				isActive    bool
				resendCount int32
				actionType  string
			)
			err := rows.Scan(&actionType, &item.SmsValue, &item.SmsCount, &isActive, &resendCount)
			if err != nil {
				return nil, err
			}
			item.ActionName = actionType
			item.IsActive = fmt.Sprintf("%v", isActive)

			if actionType == "INSUFFICIENT_BALANCE_ALERT" {
				insufficientBalanceResendCount = resendCount
			}

			datas = append(datas, &item)
		}

		return &pb.GetSmsTemplateResponse{
			Datas:                          datas,
			InsufficientBalanceResendCount: insufficientBalanceResendCount,
			SmsType:                        smsType,
		}, nil
	}

	rows, err = r.db.Query(`
		SELECT array_to_string(texts, ' '), sms_count, created_at
		FROM sms_template
		WHERE company_id = $1 AND sms_template_type = 'TEMPLATE'
	`, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datas []*pb.SmsTemplateList
	for rows.Next() {
		var item pb.SmsTemplateList
		err := rows.Scan(&item.SmsValue, &item.SmsCount, &item.ActionName)
		if err != nil {
			return nil, err
		}
		datas = append(datas, &item)
	}

	return &pb.GetSmsTemplateResponse{
		Datas:   datas,
		SmsType: smsType,
	}, nil
}

func (r SmsRepository) SendSmsDirectly(req *pb.SendSmsDirectlyRequest, companyId string) (*pb.AbsResponse, error) {
	isCyrillic := false
	for _, r := range req.SmsValue {
		if r >= 0x0400 && r <= 0x04FF {
			isCyrillic = true
			break
		}
	}

	var charsPerSms int
	if isCyrillic {
		charsPerSms = 70
	} else {
		charsPerSms = 160
	}

	smsCount := int32(len([]rune(req.SmsValue)) / charsPerSms)
	if len([]rune(req.SmsValue))%charsPerSms != 0 {
		smsCount++
	}

	var phoneNumber string
	err := r.db.QueryRow(`SELECT phone_number FROM students WHERE id = $1`, req.StudentId).Scan(&phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get student phone number: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO sms_used (id, company_id, texts, sms_count, created_by_id, created_by_name, sms_used_type, student_id)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'BY_SELF', $6)
	`, companyId, req.SmsValue, smsCount, req.CreatorId, req.CreatorName, req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert sms_used: %w", err)
	}

	if sendErr := utils.SendSMS(phoneNumber, req.SmsValue); sendErr != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to send sms: %w", sendErr)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.AbsResponse{
		Status:  200,
		Message: "ok",
	}, nil
}
