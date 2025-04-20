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
