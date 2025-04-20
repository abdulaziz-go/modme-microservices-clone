package repository

import (
	"database/sql"
	"education-service/proto/pb"
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
		SELECT comment, sms_count, sum
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
		if err := rows.Scan(&item.Comment, &item.SmsCount, &item.Sum); err != nil {
			return nil, err
		}
		datas = append(datas, &item)
	}

	return &pb.GetSmsTransactionDetailResponse{
		Datas: datas,
	}, nil
}

func (r SmsRepository) GetSmsTemplate(smsType string) (*pb.GetSmsTemplateResponse, error) {
	return nil, nil
}
