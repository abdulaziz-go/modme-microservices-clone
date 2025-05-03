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
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaymentRepository struct {
	db              *sql.DB
	educationClient *clients.EducationClient
}

func (r *PaymentRepository) AddPayment(ctx context.Context, companyId string, givenDate, sum, method, comment, studentId, actionByName, actionById, groupId string, isRefund bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	validMethods := map[string]bool{"CLICK": true, "CASH": true, "PAYME": true}
	if !validMethods[method] {
		return errors.New("invalid payment method")
	}
	amount, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		return fmt.Errorf("invalid sum amount: %v", err)
	}
	if strings.Contains(givenDate, "T") {
		givenDate = givenDate[:10]
	}
	parsedDate, err := time.Parse("2006-01-02", givenDate)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	paymentID := uuid.New()
	query := `INSERT INTO student_payments 
		(id, student_id, method, amount, given_date, comment, created_by_id, created_by_name , created_at , group_id ,payment_type, company_id , teacher_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 , $10 , $11 , $12 , $13)`
	paymentType := "ADD"
	if isRefund {
		paymentType = "REFUND"
	}
	if groupId == "" {
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), nil, paymentType, companyId, nil)
	} else {
		ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
		defer cancelFunc()
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), groupId, paymentType, companyId, r.educationClient.GetGroupDetailsByGroupId(ctx, groupId).TeacherId)
	}
	if err != nil {
		return fmt.Errorf("failed to add payment: %v", err)
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	err = r.educationClient.ChangeUserBalanceHistory(ctx, studentId, sum, givenDate, comment, "ADD", actionById, actionByName, groupId)
	if err != nil {
		return fmt.Errorf("failed to update user balance history: %v", err)
	}
	return nil
}

func (r *PaymentRepository) TakeOffPayment(ctx context.Context, companyId string, date, sum, method, comment, studentId, actionByName, actionById, groupId, studentConditionDate string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	validMethods := map[string]bool{"CLICK": true, "CASH": true, "PAYME": true}
	if !validMethods[method] {
		return errors.New("invalid payment method")
	}

	amount, err := strconv.ParseFloat(sum, 64)
	if err != nil {
		return fmt.Errorf("invalid sum amount: %v", err)
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	paymentID := uuid.New()
	query := `INSERT INTO student_payments 
		(id, student_id, method, amount, given_date, comment, payment_type, created_by_id, created_by_name , created_at, group_id , company_id , student_activation_date)
		VALUES ($1, $2, $3, $4, $5, $6, 'TAKE_OFF', $7, $8 , $9 , $10 , $11 , $12)`

	if groupId == "" {
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), nil, companyId, studentConditionDate)
	} else {
		_, err = tx.Exec(query, paymentID, studentId, method, amount, parsedDate, comment, actionById, actionByName, time.Now(), groupId, companyId, studentConditionDate)
	}
	if err != nil {
		return fmt.Errorf("failed to take off payment: %v", err)
	}

	err = r.educationClient.ChangeUserBalanceHistory(ctx, studentId, sum, date, comment, "TAKE_OFF", actionById, actionByName, groupId)
	if err != nil {
		return fmt.Errorf("failed to update user balance history: %v", err)
	}
	return nil
}

func (r *PaymentRepository) PaymentReturn(ctx context.Context, companyId, paymentId, actionByName, actionById string) (*pb.AbsResponse, error) {
	var payment struct {
		ID            string
		StudentID     string
		Method        string
		Amount        float64
		GivenDate     string
		Comment       string
		PaymentType   string
		CreatedByID   string
		CreatedByName string
		CreatedAt     string
		GroupId       string
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query := `SELECT id, student_id, method, amount, given_date, comment, payment_type, created_by_id, created_by_name, created_at , coalesce(group_id , 0)
			  FROM student_payments WHERE id = $1 and company_id=$2`
	err = tx.QueryRow(query, paymentId, companyId).Scan(
		&payment.ID,
		&payment.StudentID,
		&payment.Method,
		&payment.Amount,
		&payment.GivenDate,
		&payment.Comment,
		&payment.PaymentType,
		&payment.CreatedByID,
		&payment.CreatedByName,
		&payment.CreatedAt,
		&payment.GroupId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("payment not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment: %v", err)
	}

	deleteQuery := `DELETE FROM student_payments WHERE id = $1 and company_id=$2`
	_, err = tx.Exec(deleteQuery, paymentId, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete payment: %v", err)
	}
	if payment.PaymentType == "TAKE_OFF" {
		payment.PaymentType = "ADD"
	} else {
		payment.PaymentType = "TAKE_OFF"
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	err = r.educationClient.ChangeUserBalanceHistory(ctx, payment.StudentID, fmt.Sprintf("%.2f", payment.Amount), payment.GivenDate, payment.Comment, payment.PaymentType, actionById, actionByName, payment.GroupId)
	if err != nil {
		return nil, fmt.Errorf("failed to update user balance history: %v", err)
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "payment returned successfully",
	}, nil
}

func (r *PaymentRepository) PaymentUpdate(ctx context.Context, companyId string, paymentId string, date string, method string, userId string, comment string, debit, actionByName, actionById, groupId string) (*pb.AbsResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var paymentType string
	var oldDebit string
	query := `SELECT payment_type , amount FROM student_payments where id=$1 and company_id=$2`
	err = tx.QueryRow(query, paymentId, companyId).Scan(&paymentType, &oldDebit)
	if err != nil {
		return nil, fmt.Errorf("error checking payment existence: %v", err)
	}
	updateQuery := `UPDATE student_payments 
					SET given_date = $1, method = $2, comment = $3, amount = $4, created_by_id = $5, created_by_name = $6, group_id = $8 
					WHERE id = $7 and company_id=$9`
	if groupId == "" {
		_, err = tx.Exec(updateQuery, date, method, comment, debit, actionById, actionByName, paymentId, nil, companyId)
	} else {
		_, err = tx.Exec(updateQuery, date, method, comment, debit, actionById, actionByName, paymentId, groupId, companyId)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}

	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	_, err = r.educationClient.ChangeUserBalanceHistoryByDebit(ctx, userId, oldDebit, date, comment, paymentType, actionById, actionByName, groupId, debit)
	if err != nil {
		return nil, fmt.Errorf("failed to update user balance history: %v", err)
	}

	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "payment updated successfully",
	}, nil
}

func (r *PaymentRepository) GetMonthlyStatus(ctx context.Context, companyId string, studentId string) (*pb.GetMonthlyStatusResponse, error) {
	query := `
		SELECT 
			TO_CHAR(given_date, 'YYYY-MM') AS month, 
			SUM(CASE WHEN payment_type = 'ADD' OR payment_type='REFUND' THEN amount ELSE 0 END) AS total_add,
			SUM(CASE WHEN payment_type = 'TAKE_OFF' THEN amount ELSE 0 END) AS total_take_off
		FROM 
			student_payments 
		WHERE 
			student_id = $1 and company_id=$2
		GROUP BY 
			month 
		ORDER BY 
			month`

	rows, err := r.db.Query(query, studentId, companyId)
	if err != nil {
		return nil, fmt.Errorf("error querying monthly status: %v", err)
	}
	defer rows.Close()

	var monthlyStatus []*pb.AbsGetMonthlyStatusResponse

	for rows.Next() {
		var month string
		var totalAdd, totalTakeOff float64

		if err := rows.Scan(&month, &totalAdd, &totalTakeOff); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		balance := totalAdd - totalTakeOff

		monthlyStatus = append(monthlyStatus, &pb.AbsGetMonthlyStatusResponse{
			Month:   month,
			Balance: fmt.Sprintf("%.2f", balance),
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return &pb.GetMonthlyStatusResponse{
		MonthStatus: monthlyStatus,
	}, nil
}

func (r *PaymentRepository) GetAllPaymentsByMonth(ctx context.Context, companyId string, month string, studentId string) (*pb.GetAllPaymentsByMonthResponse, error) {
	query := `
		SELECT 
		    id,
			given_date, 
			payment_type, 
			amount, 
			comment, 
			created_by_id,
			created_by_name, 
			created_at ,
			coalesce(group_id , 0),
			method
		FROM 
			student_payments 
		WHERE 
			student_id = $1 AND 
			TO_CHAR(given_date, 'YYYY-MM') = $2 and company_id=$3
		ORDER BY 
			created_at desc`

	rows, err := r.db.Query(query, studentId, month, companyId)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %v", err)
	}
	defer rows.Close()

	var payments []*pb.AbsGetAllPaymentsByMonthResponse
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()

	for rows.Next() {
		var payment pb.AbsGetAllPaymentsByMonthResponse

		var givenDate sql.NullTime

		if err := rows.Scan(
			&payment.PaymentId,
			&givenDate,
			&payment.PaymentType,
			&payment.Amount,
			&payment.Comment,
			&payment.CreatedById,
			&payment.CreatedByName,
			&payment.CreatedAt,
			&payment.GroupId,
			&payment.Method,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if givenDate.Valid {
			payment.GivenDate = givenDate.Time.Format("2006-01-02")
		} else {
			payment.GivenDate = ""
		}

		payment.GroupName = r.educationClient.GetGroupNameById(ctx, payment.GroupId)
		payments = append(payments, &payment)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return &pb.GetAllPaymentsByMonthResponse{
		Payments: payments,
	}, nil
}

func (r *PaymentRepository) GetAllPaymentTakeOff(ctx context.Context, companyId string, from string, to string) (*pb.GetAllPaymentTakeOffResponse, error) {
	query := `
        SELECT 
            id, 
            amount,
            given_date, 
            student_id, 
            comment, 
            created_by_id, 
            created_by_name
        FROM 
            student_payments
        WHERE 
            payment_type = 'TAKE_OFF' 
            AND created_by_id != '00000000-0000-0000-0000-000000000000'
            AND given_date BETWEEN $1 AND $2 and company_id=$3
    `

	rows, err := r.db.Query(query, from, to, companyId)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %w", err)
	}
	defer rows.Close()

	response := &pb.GetAllPaymentTakeOffResponse{}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	for rows.Next() {
		var payment pb.AbsPaymentTakeOff

		err := rows.Scan(
			&payment.PaymentId,
			&payment.Sum,
			&payment.GivenDate,
			&payment.StudentId,
			&payment.Comment,
			&payment.CreatorId,
			&payment.CreatorName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		name, _, _, err := r.educationClient.GetStudentById(ctx, payment.StudentId)
		if err != nil {
			payment.StudentName = "error while getting this student name"
		}
		payment.StudentName = name
		response.Pennies = append(response.Pennies, &payment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return response, nil
}

func (r *PaymentRepository) GetAllPaymentTakeOffChart(ctx context.Context, companyId string, from string, to string) (*pb.GetAllPaymentTakeOffChartResponse, error) {
	query := `
        SELECT 
            given_date, 
            SUM(amount)
        FROM 
            student_payments
        WHERE 
            payment_type = 'TAKE_OFF'
            AND created_by_id != '00000000-0000-0000-0000-000000000000'
            AND given_date BETWEEN $1 AND $2 and company_id=$3
        GROUP BY 
            given_date
        ORDER BY 
            given_date;
    `

	rows, err := r.db.Query(query, from, to, companyId)
	if err != nil {
		return nil, fmt.Errorf("error querying payment chart data: %w", err)
	}
	defer rows.Close()

	response := &pb.GetAllPaymentTakeOffChartResponse{}

	for rows.Next() {
		var chartEntry pb.AbsTakeOfChartResponse

		err := rows.Scan(
			&chartEntry.YearMonth,
			&chartEntry.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chart data row: %w", err)
		}

		response.ChartResponse = append(response.ChartResponse, &chartEntry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return response, nil
}
func (r *PaymentRepository) GetAllStudentPayments(
	ctx context.Context,
	companyId string,
	from string,
	to string,
	filters []*pb.Filters,
	sorts []*pb.SortBy,
	pageRequest *pb.PageRequest,
) (*pb.GetAllStudentPaymentsResponse, error) {
	query := `
SELECT 
    student_id,
    method,
    amount,
    given_date,
    comment,
    created_by_id,
    created_by_name
FROM student_payments
WHERE given_date BETWEEN $1 AND $2
  AND payment_type = 'ADD'
  AND company_id = $3
`

	args := []interface{}{from, to, companyId}
	argIndex := 4

	// Add filters
	allowedFields := map[string]bool{
		"method":        true,
		"group_id":      true,
		"student_id":    true,
		"created_by_id": true,
		"amount":        true,
		"teacher_id":    true,
	}

	for _, f := range filters {
		if f.Field == "" || f.Value == "" || !allowedFields[f.Field] {
			continue
		}
		operator := getSQLOperator(f.Type)
		if operator == "" {
			continue
		}
		query += fmt.Sprintf(" AND %s %s $%d", f.Field, operator, argIndex)
		args = append(args, f.Value)
		argIndex++
	}

	// Sorting
	if len(sorts) > 0 {
		query += " ORDER BY "
		orderClauses := []string{}
		for _, s := range sorts {
			order := "ASC"
			if strings.ToUpper(s.Type) == "DESC" {
				order = "DESC"
			}
			orderClauses = append(orderClauses, fmt.Sprintf("%s %s", s.Field, order))
		}
		query += strings.Join(orderClauses, ", ")
	} else {
		query += " ORDER BY created_at DESC"
	}

	// Pagination
	if pageRequest != nil && pageRequest.Size > 0 {
		offset := (pageRequest.Page - 1) * pageRequest.Size
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, pageRequest.Size, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ctx, cancel := utils.NewTimoutContext(ctx, companyId)
	defer cancel()

	resp := pb.GetAllStudentPaymentsResponse{}

	for rows.Next() {
		el := pb.AbsStudentPayments{}
		err := rows.Scan(&el.StudentId, &el.Method, &el.Amount, &el.GivenDate, &el.Comment, &el.CreatorId, &el.CreatorName)
		if err != nil {
			return nil, err
		}
		name, _, _, _ := r.educationClient.GetStudentById(ctx, el.StudentId)
		el.StudentName = name
		resp.Payments = append(resp.Payments, &el)
	}

	return &resp, nil
}

func (r *PaymentRepository) GetAllStudentPaymentsChart(
	ctx context.Context,
	companyId string,
	from string,
	to string,
	filters []*pb.Filters,
	sorts []*pb.SortBy,
	_ *pb.PageRequest,
) (*pb.GetAllStudentPaymentsChartResponse, error) {
	var (
		cash, payme, click float64
	)

	resp := pb.GetAllStudentPaymentsChartResponse{}

	baseConditions := []string{
		"payment_type != 'TAKE_OFF'",
		"given_date BETWEEN $1 AND $2",
		"company_id = $3",
	}

	args := []interface{}{from, to, companyId}
	argIndex := 4

	allowedFields := map[string]bool{
		"method":        true,
		"group_id":      true,
		"student_id":    true,
		"created_by_id": true,
		"amount":        true,
	}

	for _, f := range filters {
		if f.Field == "" || f.Value == "" || !allowedFields[f.Field] {
			continue
		}
		operator := getSQLOperator(f.Type)
		if operator == "" {
			continue
		}
		baseConditions = append(baseConditions, fmt.Sprintf("%s %s $%d", f.Field, operator, argIndex))
		args = append(args, f.Value)
		argIndex++
	}

	conditionStr := strings.Join(baseConditions, " AND ")

	summaryQuery := fmt.Sprintf(`
SELECT
  COALESCE((SELECT SUM(amount) FROM student_payments WHERE method = 'CASH' AND %s), 0),
  COALESCE((SELECT SUM(amount) FROM student_payments WHERE method = 'PAYME' AND %s), 0),
  COALESCE((SELECT SUM(amount) FROM student_payments WHERE method = 'CLICK' AND %s), 0)
`, conditionStr, conditionStr, conditionStr)

	err := r.db.QueryRow(summaryQuery, args...).Scan(&cash, &payme, &click)
	if err != nil {
		return nil, fmt.Errorf("error scanning revenue summary: %w", err)
	}

	resp.TotalRevenue = strconv.FormatFloat(cash+payme+click, 'f', 2, 64)
	resp.Cash = strconv.FormatFloat(cash, 'f', 2, 64)
	resp.Click = strconv.FormatFloat(click, 'f', 2, 64)
	resp.Payme = strconv.FormatFloat(payme, 'f', 2, 64)

	chartQuery := `
SELECT 
  given_date, 
  SUM(amount)
FROM 
  student_payments
WHERE 
  payment_type = 'ADD' 
  AND given_date BETWEEN $1 AND $2 
  AND company_id = $3
`
	chartArgs := []interface{}{from, to, companyId}
	chartIndex := 4

	for _, f := range filters {
		if f.Field == "" || f.Value == "" || !allowedFields[f.Field] {
			continue
		}
		operator := getSQLOperator(f.Type)
		if operator == "" {
			continue
		}
		chartQuery += fmt.Sprintf(" AND %s %s $%d", f.Field, operator, chartIndex)
		chartArgs = append(chartArgs, f.Value)
		chartIndex++
	}

	chartQuery += `
GROUP BY given_date
`

	// Sorting
	if len(sorts) > 0 {
		chartQuery += " ORDER BY "
		orderParts := []string{}
		for _, s := range sorts {
			order := "ASC"
			if strings.ToUpper(s.Type) == "DESC" {
				order = "DESC"
			}
			switch s.Field {
			case "given_date":
				orderParts = append(orderParts, fmt.Sprintf("%s %s", s.Field, order))
			}
		}
		if len(orderParts) > 0 {
			chartQuery += strings.Join(orderParts, ", ")
		} else {
			chartQuery += "given_date"
		}
	} else {
		chartQuery += " ORDER BY given_date"
	}

	rows, err := r.db.Query(chartQuery, chartArgs...)
	if err != nil {
		return nil, fmt.Errorf("error querying chart data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chartEntry pb.AbsTakeOfChartResponse
		if err := rows.Scan(&chartEntry.YearMonth, &chartEntry.Amount); err != nil {
			return nil, fmt.Errorf("error scanning chart row: %w", err)
		}
		resp.PaymentsChart = append(resp.PaymentsChart, &chartEntry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return &resp, nil
}

func getSQLOperator(op string) string {
	switch op {
	case "=", "==":
		return "="
	case "!=", "<>":
		return "!="
	case "<":
		return "<"
	case "<=":
		return "<="
	case ">":
		return ">"
	case ">=":
		return ">="
	default:
		return ""
	}
}

func (r *PaymentRepository) GetAllDebtsInformation(ctx context.Context, companyId string, from, to string, amountFrom, amountTo int64, page, size int32) (*pb.GetAllDebtsInformationResponse, error) {
	var (
		query string
		rows  *sql.Rows
		err   error
	)
	if size <= 0 {
		return nil, fmt.Errorf("invalid page size: must be greater than zero")
	}
	offset := (page - 1) * size

	if amountFrom != 0 || amountTo != 0 {
		query = `
                        SELECT student_id AS debtor_id,
                                COALESCE(SUM(CASE WHEN payment_type = 'ADD' OR payment_type = 'REFUND' THEN amount ELSE 0 END), 0) -
                                COALESCE(SUM(CASE WHEN payment_type = 'TAKE_OFF' THEN amount ELSE 0 END), 0) AS total_on_period
                        FROM student_payments
                        WHERE given_date BETWEEN $1 AND $2 AND company_id = $5
                        GROUP BY student_id
                        HAVING COALESCE(SUM(CASE WHEN payment_type = 'ADD' OR payment_type = 'REFUND' THEN amount ELSE 0 END), 0) -
                               COALESCE(SUM(CASE WHEN payment_type = 'TAKE_OFF' THEN amount ELSE 0 END), 0) BETWEEN $6 AND $7
                        LIMIT $3 OFFSET $4;
                `
		rows, err = r.db.Query(query, from, to, size, offset, companyId, amountFrom, amountTo)
		if err != nil {
			return nil, err
		}
	} else {
		query = `
                        SELECT student_id AS debtor_id,
                                COALESCE(SUM(CASE WHEN payment_type = 'ADD' OR payment_type = 'REFUND' THEN amount ELSE 0 END), 0) -
                                COALESCE(SUM(CASE WHEN payment_type = 'TAKE_OFF' THEN amount ELSE 0 END), 0) AS total_on_period
                        FROM student_payments
                        WHERE given_date BETWEEN $1 AND $2 AND company_id = $5
                        GROUP BY student_id
                        LIMIT $3 OFFSET $4;
                `
		rows, err = r.db.Query(query, from, to, size, offset, companyId)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	var debts []*pb.AbsDebtsInformation
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	for rows.Next() {
		var debt pb.AbsDebtsInformation
		if err := rows.Scan(&debt.DebtorId, &debt.TotalOnPeriod); err != nil {
			return nil, err
		}
		name, phoneNumber, balance, err := r.educationClient.GetStudentById(ctx, debt.DebtorId)
		if err != nil {
			name = "unknown name"
			phoneNumber = "unknown phoneNumber"
			balance = 0
		}
		if balance >= 0 {
			continue
		}
		debt.DebtorName = name
		debt.PhoneNumber = phoneNumber
		resp, _ := r.educationClient.GetGroupsAndCommentsByStudentId(ctx, debt.DebtorId)
		debt.Comments = resp.Comments
		debt.Groups = resp.Groups
		debt.Balance = strconv.FormatFloat(balance, 'f', 2, 64)
		debts = append(debts, &debt)
	}

	var totalRecords int32
	var countQuery string
	if amountFrom != 0 || amountTo != 0 {
		countQuery = `
    SELECT COUNT(*) FROM (
        SELECT student_id
        FROM student_payments
        WHERE given_date BETWEEN $1 AND $2 AND company_id = $3
        GROUP BY student_id
        HAVING COALESCE(SUM(CASE WHEN payment_type = 'ADD' OR payment_type = 'REFUND' THEN amount ELSE 0 END), 0) -
               COALESCE(SUM(CASE WHEN payment_type = 'TAKE_OFF' THEN amount ELSE 0 END), 0) BETWEEN $4 AND $5
    ) AS filtered_students`

		if err := r.db.QueryRow(countQuery, from, to, companyId, amountFrom, amountTo).Scan(&totalRecords); err != nil {
			return nil, err
		}
	} else {
		countQuery = `SELECT COUNT(*) FROM (
    SELECT student_id
    FROM student_payments
    WHERE given_date BETWEEN $1 AND $2 AND company_id = $3
    GROUP BY student_id
) AS filtered_students
`
		if err := r.db.QueryRow(countQuery, from, to, companyId).Scan(&totalRecords); err != nil {
			return nil, err
		}
	}
	if totalRecords == 0 {
		totalRecords = 1
	}
	totalPageCount := (totalRecords + size - 1) / size

	return &pb.GetAllDebtsInformationResponse{
		TotalPageCount: totalPageCount,
		Debts:          debts,
	}, nil
}

func (r *PaymentRepository) GetCommonFinanceInformation(ctx context.Context, companyId string) (*pb.GetCommonInformationResponse, error) {
	response := new(pb.GetCommonInformationResponse)

	var payInCurrentMonth int32

	err := r.db.QueryRow(`
        SELECT COUNT(id)
        FROM student_payments
        WHERE payment_type = 'ADD' and company_id=$1
        AND EXTRACT(MONTH FROM given_date) = EXTRACT(MONTH FROM CURRENT_DATE) 
        AND EXTRACT(YEAR FROM given_date) = EXTRACT(YEAR FROM CURRENT_DATE);
    `, companyId).Scan(&payInCurrentMonth)
	if err != nil {
		payInCurrentMonth = 0
	}

	response.DebtorsCount = 0
	response.PayInCurrentMonth = payInCurrentMonth

	return response, nil
}

func (r *PaymentRepository) GetIncomeChart(ctx context.Context, companyId string, from string, to string) (*pb.GetIncomeChartResponse, error) {
	startDate, err := time.Parse("200601", from)
	if err != nil {
		return nil, fmt.Errorf("invalid from date format: %v", err)
	}
	endDate, err := time.Parse("200601", to)
	if err != nil {
		return nil, fmt.Errorf("invalid to date format: %v", err)
	}

	query := `
        SELECT 
            TO_CHAR(given_date, 'YYYYMM') AS specific_month,
            SUM(
                CASE 
                    WHEN payment_type IN ('ADD', 'TAKE_OFF') THEN amount
                    WHEN payment_type = 'REFUND' THEN -amount
                    ELSE 0
                END
            ) AS balance
        FROM student_payments
        WHERE given_date BETWEEN $1 AND $2 AND company_id=$3
        GROUP BY TO_CHAR(given_date, 'YYYYMM')
        ORDER BY specific_month;
    `

	rows, err := r.db.Query(query, startDate, endDate, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monthlyBalances := make(map[string]float64)

	for rows.Next() {
		var month string
		var balance float64
		if err := rows.Scan(&month, &balance); err != nil {
			return nil, err
		}
		monthlyBalances[month] = balance
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var response pb.GetIncomeChartResponse
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 1, 0) {
		month := d.Format("200601")
		balance := monthlyBalances[month] // zero if not found

		response.Response = append(response.Response, &pb.AbsIncomeChart{
			SpecificMonth: month,
			Balance:       fmt.Sprintf("%.2f", balance),
		})
	}
	return &response, nil
}

func NewPaymentRepository(db *sql.DB, client *clients.EducationClient) *PaymentRepository {
	return &PaymentRepository{db: db, educationClient: client}
}
