package repository

import (
	"context"
	"database/sql"
	"errors"
	"finance-service/internal/clients"
	"finance-service/internal/utils"
	"finance-service/proto/pb"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type DiscountRepository struct {
	db            *sql.DB
	studentClient *clients.EducationClient
	paymentRepo   *PaymentRepository
}

func (r *DiscountRepository) CreateDiscount(ctx context.Context, companyId, groupId string, studentId string, discountPrice, comment, startDate, endDate string, withTeacher bool) error {
	var checker bool
	tx, err := r.db.Begin()
	if err != nil {
		return status.Errorf(codes.Aborted, "error while creating transaction %v", err)
	}

	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM student_discount WHERE group_id=$1 AND student_id=$2 AND company_id=$3)`, groupId, studentId, companyId).Scan(&checker)
	if err != nil {
		return status.Errorf(codes.Internal, "Error checking existing discount: %v", err)
	}
	if checker {
		return status.Errorf(codes.AlreadyExists, "Discount already exists")
	}

	_, err = tx.Exec(`INSERT INTO student_discount (student_id, discount, group_id, comment, start_at, end_at, withteacher , company_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7 , $8)`, studentId, discountPrice, groupId, comment, startDate, endDate, withTeacher, companyId)
	if err != nil {
		tx.Rollback()
		return status.Errorf(codes.Internal, "Status discount insert error: %v", err)
	}

	_, err = tx.Exec(`INSERT INTO student_discount_history (id, student_id, group_id, start_at, end_at, withteacher, comment, action , discount, company_id)
		VALUES ($1, $2, $3, $4, $5 , $6 , $7, $8, $9, $10)`, uuid.New(), studentId, groupId, startDate, endDate, withTeacher, comment, "CREATE", discountPrice, companyId)
	if err != nil {
		tx.Rollback()
		return status.Errorf(codes.Internal, "Error inserting into student history: %v", err)
	}

	rows, err := tx.Query(`
        SELECT id, student_id, method, amount, given_date, comment, created_at, payment_type, created_by_id, created_by_name, group_id , student_activation_date
        FROM student_payments
        WHERE student_id = $1 AND group_id = $2 AND given_date BETWEEN $3 AND $4 AND payment_type = 'TAKE_OFF' AND company_id=$5 AND student_activation_date is not null`,
		studentId, groupId, startDate, endDate, companyId)
	if err != nil {
		tx.Rollback()
		return status.Errorf(codes.Internal, "Error fetching student payments: %v", err)
	}
	defer rows.Close()
 
	var payments []struct {
		ID                    string
		StudentID             string
		Method                string
		Amount                float64
		GivenDate             string
		Comment               string
		CreatedAt             string
		PaymentType           string
		CreatedByID           string
		CreatedByName         string
		GroupID               int64
		StudentActivationDate string
	}
	for rows.Next() {
		var payment struct {
			ID                    string
			StudentID             string
			Method                string
			Amount                float64
			GivenDate             string
			Comment               string
			CreatedAt             string
			PaymentType           string
			CreatedByID           string
			CreatedByName         string
			GroupID               int64
			StudentActivationDate string
		}
		if err := rows.Scan(&payment.ID, &payment.StudentID, &payment.Method, &payment.Amount, &payment.GivenDate, &payment.Comment, &payment.CreatedAt, &payment.PaymentType, &payment.CreatedByID, &payment.CreatedByName, &payment.GroupID, &payment.StudentActivationDate); err != nil {
			tx.Rollback()
			return status.Errorf(codes.Internal, "Error scanning student payments: %v", err)
		}
		payments = append(payments, payment)
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	for _, payment := range payments {
		discountAmount, err := r.studentClient.CalculateDiscountSumma(ctx, groupId, studentId, discountPrice, startDate, endDate, payment.GivenDate, payment.StudentActivationDate)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = r.paymentRepo.AddPayment(ctx, companyId, payment.GivenDate, discountAmount, "CASH", "Studentga ushbu tolov amalga oshirilgan kunlar oralig'ida chegirma kiritildi va studentning qolgan puli qaytarib berildi.", studentId, payment.CreatedByName, payment.CreatedByID, groupId, true)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return status.Errorf(codes.Aborted, "Transaction commit failed: %v", err)
	}

	return nil
}

func (r *DiscountRepository) DeleteDiscount(companyId, groupId string, studentId string) error {
	var discount pb.AbsDiscountRequest
	var createdAt string
	err := r.db.QueryRow(`SELECT student_id, discount, group_id, comment, start_at, end_at, withteacher, created_at FROM student_discount where student_id=$1 and group_id=$2 and company_id=$3 `, studentId, groupId, companyId).Scan(
		&discount.StudentId,
		&discount.DiscountPrice,
		&discount.GroupId,
		&discount.Comment,
		&discount.StartDate,
		&discount.EndDate,
		&discount.WithTeacher,
		&createdAt,
	)
	if err != nil {
		return status.Errorf(codes.NotFound, err.Error())
	}
	_, err = r.db.Exec(`DELETE FROM student_discount where group_id=$1 and student_id=$2 and company_id=$3`, groupId, studentId, companyId)
	if err != nil {
		return status.Errorf(codes.Internal, "Error while deleting disount %f", err)
	}
	_, err = r.db.Exec(`INSERT INTO student_discount_history (id, student_id, group_id, start_at, end_at, withteacher, comment, action , discount , company_id)
		VALUES ($1, $2, $3, $4, $5 , $6 , $7, $8 , $9 , $10)`, uuid.New(), studentId, groupId, discount.StartDate, discount.EndDate, discount.WithTeacher, discount.Comment, "DELETE", discount.DiscountPrice, companyId)
	if err != nil {
		return status.Errorf(codes.Internal, "Error inserting into student history: %v", err)
	}
	return nil
}
func (r *DiscountRepository) GetAllDiscountByGroup(ctx context.Context, companyId, groupId string) (*pb.GetInformationDiscountResponse, error) {
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	resp, err := r.studentClient.GetStudentsByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}
	result := pb.GetInformationDiscountResponse{}
	var res []*pb.AbsStudentDiscount
	students := resp.Students
	for _, el := range students {
		var element pb.AbsStudentDiscount
		element.StudentId = el.Id
		element.StudentName = el.Name
		element.StudentPhoneNumber = el.PhoneNumber
		r.db.QueryRow(`SELECT  discount, comment, created_at , start_at , end_at , withteacher FROM student_discount WHERE group_id=$1 and student_id=$2 and company_id=$3`, groupId, el.Id, companyId).Scan(&element.Discount, &element.Cause, &element.CreatedAt, &element.StartAt, &element.EndAt, &element.WithTeacher)
		res = append(res, &element)
	}
	result.Discounts = res
	return &result, nil
}

func (r *DiscountRepository) GetHistoryDiscount(ctx context.Context, companyId, id string) (*pb.GetHistoryDiscountResponse, error) {
	query := `
		SELECT group_id, student_id,discount, comment, start_at, end_at, withTeacher, action, created_at
		FROM student_discount_history
		WHERE group_id = $1
	`

	var discounts []*pb.AbsHistoryDiscount

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	for rows.Next() {
		var discount pb.AbsHistoryDiscount
		var startAt, endAt, createdAt time.Time

		if err := rows.Scan(
			&discount.GroupId,
			&discount.StudentId,
			&discount.DiscountPrice,
			&discount.Comment,
			&startAt,
			&endAt,
			&discount.WithTeacher,
			&discount.Action,
			&createdAt,
		); err != nil {
			return nil, err
		}

		name, _, _, _ := r.studentClient.GetStudentById(ctx, discount.StudentId)
		discount.StudentName = name
		discount.StartDate = startAt.Format("2006-01-02")
		discount.EndDate = endAt.Format("2006-01-02")
		discount.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		discounts = append(discounts, &discount)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetHistoryDiscountResponse{Discounts: discounts}, nil
}

func (r *DiscountRepository) GetDiscountByStudentId(companyId, studentId, groupId string) (*pb.GetDiscountByStudentIdResponse, error) {
	var discount float64
	var startAt, endAt string
	var withTeacher bool

	err := r.db.QueryRow(`SELECT discount, start_at, end_at , withteacher FROM student_discount WHERE student_id=$1 AND group_id=$2 and company_id=$3 order by created_at desc limit 1`, studentId, groupId, companyId).Scan(&discount, &startAt, &endAt, &withTeacher)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no discount found for the given student and group")
		}
		return nil, fmt.Errorf("failed to query database: %v", err)
	}

	startTime, err := time.Parse("2006-01-02T15:04:05Z", startAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("failed to parse start_at: %v", err)
	}

	endTime, err := time.Parse("2006-01-02T15:04:05Z", endAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_at: %v", err)
	}

	now := time.Now()
	if now.After(startTime) && now.Before(endTime) {
		response := &pb.GetDiscountByStudentIdResponse{}
		if withTeacher {
			response.Amount = fmt.Sprintf("%.2f", discount)
			response.IsHave = true
			response.DiscountOwner = "TEACHER"
		} else {
			response.Amount = fmt.Sprintf("%.2f", discount)
			response.IsHave = true
			response.DiscountOwner = "CENTER"
		}
		return response, nil
	}

	return nil, fmt.Errorf("current time is not within the discount period")
}

func NewDiscountRepository(db *sql.DB, studentClient *clients.EducationClient, paymentRepo *PaymentRepository) *DiscountRepository {
	return &DiscountRepository{db: db, studentClient: studentClient, paymentRepo: paymentRepo}
}
