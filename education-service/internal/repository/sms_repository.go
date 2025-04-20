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
	baseQuery := `SELECT sms_count, array_to_string(texts, ' '), created_by_name, created_at FROM sms_used WHERE company_id = $1`
	args := []interface{}{companyId}

	if len(studentId) != 0 {
		baseQuery += " AND student_id = $2"
		args = append(args, studentId)
	}

	baseQuery += " ORDER BY created_at DESC LIMIT $3 OFFSET $4"
	args = append(args, pageReq.Size, offset)

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
