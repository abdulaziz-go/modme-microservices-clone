package repository

import (
	"context"
	"database/sql"
	"education-service/internal/clients"
	"education-service/internal/models"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StudentRepository struct {
	db                *sql.DB
	userClient        *clients.UserClient
	financeClient     *clients.FinanceClient
	financeClientChan chan *clients.FinanceClient
}

func NewStudentRepository(db *sql.DB, userClient *clients.UserClient, financeClientChan chan *clients.FinanceClient) *StudentRepository {
	return &StudentRepository{db: db, userClient: userClient, financeClientChan: financeClientChan}
}

func (r *StudentRepository) ensureFinanceClient() error {
	if r.financeClient == nil {
		select {
		case client := <-r.financeClientChan:
			r.financeClient = client
		case <-time.After(5 * time.Second):
			return fmt.Errorf("failed to initialize GroupClient within timeout")
		}
	}
	return nil
}
func (r *StudentRepository) GetAllStudent(ctx context.Context, companyId string, condition string, page string, size string, courseId string, nameOrSurname string, trainingDateStart, trainingDateEnd string) (*pb.GetAllStudentResponse, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		return nil, fmt.Errorf("invalid page value: %v", err)
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil || sizeInt < 1 {
		return nil, fmt.Errorf("invalid size value: %v", err)
	}
	offset := (pageInt - 1) * sizeInt

	whereConditions := []string{"condition = $1", "company_id = $2"}
	args := []interface{}{condition, companyId}
	paramCount := 2

	if nameOrSurname != "" && nameOrSurname[0] != '+' {
		paramCount++
		whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", paramCount))
		args = append(args, "%"+nameOrSurname+"%")
	}
	if nameOrSurname != "" && nameOrSurname[0] == '+' {
		paramCount++
		whereConditions = append(whereConditions, fmt.Sprintf("phone ILIKE $%d", paramCount))
		args = append(args, "%"+nameOrSurname+"%")
	}

	whereClause := strings.Join(whereConditions, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students WHERE %s", whereClause)
	var totalCount int32
	err = r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(sizeInt)))

	studentQuery := fmt.Sprintf(`
    SELECT id, name, gender, date_of_birth, phone, address, passport_id, additional_contact, 
           balance, condition, telegram_username, created_at
    FROM students
    WHERE %s
    ORDER BY created_at desc 
    LIMIT $%d OFFSET $%d`, whereClause, paramCount+1, paramCount+2)

	args = append(args, sizeInt, offset)

	studentRows, err := r.db.Query(studentQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute student query: %v", err)
	}
	defer studentRows.Close()

	var students []*pb.GetGroupsAbsForStudent
	for studentRows.Next() {
		var student pb.GetGroupsAbsForStudent
		err := studentRows.Scan(
			&student.Id, &student.Name, &student.Gender, &student.DateOfBirth, &student.Phone,
			&student.Address, &student.PassportId, &student.AdditionalContact, &student.Balance,
			&student.Condition, &student.TelegramUsername, &student.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student row: %v", err)
		}
		students = append(students, &student)
	}

	if len(students) == 0 {
		return &pb.GetAllStudentResponse{Response: []*pb.GetGroupsAbsForStudent{}, TotalCount: totalCount}, nil
	}

	studentIDs := make([]string, len(students))
	for i, student := range students {
		studentIDs[i] = student.Id
	}

	groupQueryWhere := "s.id = ANY($1) AND c.company_id = $2"
	groupQueryArgs := []interface{}{pq.Array(studentIDs), companyId}
	groupParamCount := 2

	if courseId != "" {
		groupParamCount++
		groupQueryWhere += fmt.Sprintf(" AND c.id = $%d", groupParamCount)
		groupQueryArgs = append(groupQueryArgs, courseId)
	}

	if trainingDateStart != "" {
		groupParamCount++
		groupQueryWhere += fmt.Sprintf(" AND g.start_date >= $%d", groupParamCount)
		groupQueryArgs = append(groupQueryArgs, trainingDateStart)
	}

	if trainingDateEnd != "" {
		groupParamCount++
		groupQueryWhere += fmt.Sprintf(" AND g.start_date <= $%d", groupParamCount)
		groupQueryArgs = append(groupQueryArgs, trainingDateEnd)
	}

	groupQuery := fmt.Sprintf(`
    SELECT 
        s.id AS student_id,
        g.id AS group_id, g.name AS group_name, g.start_date, g.end_date, g.date_type, g.days, g.start_time,
        c.id AS course_id, c.title AS course_title, c.duration_lesson, c.course_duration, c.price,
        g.teacher_id,
        gs.condition AS student_group_condition, g.room_id, gs.last_specific_date AS student_activated_at
    FROM students s
    JOIN group_students gs ON s.id = gs.student_id
    JOIN groups g ON gs.group_id = g.id AND g.is_archived = 'false'
    JOIN courses c ON g.course_id = c.id
    WHERE %s
    ORDER BY s.id, g.id`, groupQueryWhere)

	groupRows, err := r.db.Query(groupQuery, groupQueryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute group query: %v", err)
	}
	defer groupRows.Close()

	studentMap := make(map[string]*pb.GetGroupsAbsForStudent)
	for _, student := range students {
		studentMap[student.Id] = student
	}
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()

	for groupRows.Next() {
		var studentID string
		var group pb.GroupGetAllStudentAbs
		var course pb.AbsCourse
		var teacherId string
		err := groupRows.Scan(
			&studentID,
			&group.Id, &group.Name, &group.GroupStartDate, &group.GroupEndDate, &group.Type, pq.Array(&group.Days), &group.LessonStartTime,
			&course.Id, &course.Name, &course.LessonDuration, &course.CourseDuration, &course.Price,
			&teacherId, &group.StudentCondition, &group.RoomId, &group.StudentActivatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group row: %v", err)
		}

		name, err := r.userClient.GetTeacherById(ctx, teacherId)
		if err != nil {
			return nil, err
		}

		group.TeacherName = name
		group.Course = &course
		student := studentMap[studentID]
		student.Groups = append(student.Groups, &group)
	}

	return &pb.GetAllStudentResponse{
		Response:   students,
		TotalCount: totalPages,
	}, nil
}
func (r *StudentRepository) CreateStudent(companyId string, createdBy string, phoneNumber string, name string, groupId string, address string, additionalContact string, dateFrom string, birthDate string, gender bool, passportId string, telegramUsername string) error {
	studentId := uuid.New()
	_, err := r.db.Exec(`INSERT INTO students(id, name, phone, date_of_birth, gender, telegram_username, passport_id, additional_contact, address , company_id) values ($1, $2,$3,$4,$5,$6,$7,$8,$9 , $10)`, studentId, name, phoneNumber, birthDate, gender, telegramUsername, passportId, additionalContact, address, companyId)
	if err != nil {
		return err
	}
	if groupId != "" && dateFrom != "" && createdBy != "" {
		_, err = r.db.Exec(`INSERT INTO group_students(id, group_id, student_id, created_by , company_id) values ($1 ,$2 ,$3 ,$4 , $5)`, uuid.New(), groupId, studentId, createdBy, companyId)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *StudentRepository) UpdateStudent(companyId string, studentId string, number string, name string, address string, additionalContact string, birth string, gender bool, passportId string) error {
	_, err := r.db.Exec(`UPDATE students SET phone =$1, name=$2, address =$3, additional_contact =$4, date_of_birth =$5, gender =$6, passport_id=$7 where id=$8 and company_id=$9`, number, name, address, additionalContact, birth, gender, passportId, studentId, companyId)
	if err != nil {
		return err
	}
	return nil
}
func (r *StudentRepository) DeleteStudent(ctx context.Context, companyId string, studentId string, returnMoney bool, actionById, actionByName string) error {
	var cond string
	if err := r.db.QueryRow(`select condition from students where id = $1 and company_id=$2`, studentId, companyId).Scan(&cond); err != nil {
		return err
	}

	if cond == "ACTIVE" {
		groupsQuery := `SELECT  group_id FROM group_students where condition='ACTIVE' and student_id=$1 and company_id=$2`
		rows, err := r.db.Query(groupsQuery, studentId, companyId)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var groupId string
			err = rows.Scan(&groupId)
			if err != nil {
				return err
			}
			_, _ = r.ChangeConditionStudent(ctx, companyId, studentId, groupId, "DELETE", returnMoney, time.Now().Format("2006-01-02"), actionById, actionByName, "")
		}
		_, err = r.db.Exec(`UPDATE students SET condition='ARCHIVED' where id=$1 and company_id=$2`, studentId, companyId)
		if err != nil {
			return err
		}
	} else {
		_, err := r.db.Exec(`UPDATE students SET condition='ACTIVE' where id=$1 and company_id=$2`, studentId, companyId)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *StudentRepository) AddToGroup(ctx context.Context, companyId string, groupId string, studentIds []string, createdDate, createdBy string) error {
	var checker bool
	query := `INSERT INTO group_students(id, group_id, student_id, condition, last_specific_date, created_by , company_id) values ($1 ,$2 ,$3 ,$4 , $5 , $6 , $7)`
	queryForChecking := `SELECT exists(SELECT 1 FROM students where condition = 'ARCHIVED' and id=$1 and company_id=$2)`
	queryGroupChecking := `SELECT exists(SELECT 1 FROM groups where id=$1 and is_archived=true and company_id=$2)`
	err := r.db.QueryRow(queryGroupChecking, groupId, companyId).Scan(&checker)
	if err != nil || checker {
		return errors.New(fmt.Sprintf("forbidden (archived group action error) id=%s", groupId))
	}
	for _, data := range studentIds {
		err := r.db.QueryRow(queryForChecking, data, companyId).Scan(&checker)
		if err != nil || checker {
			return errors.New(fmt.Sprintf("forbidden (archived student action error) id=%s", data))
		}
		_, err = r.db.Exec(query, uuid.New(), groupId, data, "FREEZE", createdDate, createdBy, companyId)
		if err != nil {
			continue
		}
		r.SendSmsReasonAddToGroup(ctx, data, groupId, companyId)
	}
	return nil
}

func (r *StudentRepository) GetStudentById(ctx context.Context, companyId string, id string) (*pb.GetStudentByIdResponse, error) {
	if err := r.ensureFinanceClient(); err != nil {
		return nil, err
	}

	var result pb.GetStudentByIdResponse

	err := r.db.QueryRow(`SELECT id, name, gender, date_of_birth, phone, balance, created_at , condition , additional_contact
                          FROM students WHERE id = $1 and company_id=$2`, id, companyId).
		Scan(&result.Id, &result.Name, &result.Gender, &result.DateOfBirth, &result.Phone, &result.Balance, &result.CreatedAt, &result.Condition, &result.AdditionalContact)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`
        SELECT gs.created_at, gs.group_id, gs.condition, gs.last_specific_date, 
               g.name, g.date_type, g.days, g.start_time, g.start_date, g.end_date,
               r.id, r.capacity, r.title, 
               c.id, c.title, c.duration_lesson, c.course_duration, c.price, c.description , g.teacher_id
        FROM group_students gs
        JOIN groups g ON g.id = gs.group_id and g.is_archived='false'
        JOIN rooms r ON g.room_id = r.id
        JOIN courses c ON c.id = g.course_id
        WHERE gs.student_id = $1 and c.company_id=$2`, id, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	for rows.Next() {
		var groupStudent pb.GetGroupStudent
		var room pb.AbsRoom
		var course pb.AbsCourse
		var teacherId string

		err := rows.Scan(&groupStudent.StudentAddedAt, &groupStudent.Id, &groupStudent.StudentCondition, &groupStudent.StudentActivatedAt,
			&groupStudent.Name, &groupStudent.DateType, pq.Array(&groupStudent.Days), &groupStudent.LessonStartTime,
			&groupStudent.GroupStartDate, &groupStudent.GroupEndDate,
			&room.Id, &room.Capacity, &room.Name,
			&course.Id, &course.Name, &course.LessonDuration, &course.CourseDuration, &course.Price, &course.Description, &teacherId)
		if err != nil {
			return nil, err
		}

		discount, _ := r.financeClient.GetDiscountByStudentId(ctx, result.Id, groupStudent.Id)
		groupStudent.PriceForStudent = course.Price
		if discount != nil {
			groupStudent.PriceForStudent = course.Price - *discount
		}
		groupStudent.Room = &room
		groupStudent.Course = &course
		name, err := r.userClient.GetTeacherById(ctx, teacherId)
		if err != nil {
			return nil, err
		}
		groupStudent.TeacherName = name
		result.Groups = append(result.Groups, &groupStudent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}
func (r *StudentRepository) GetNoteByStudent(companyId string, id string) (*pb.GetNotesByStudent, error) {
	rows, err := r.db.Query(`SELECT id, comment, created_at FROM student_note where student_id=$1 and company_id=$2`, id, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []*pb.AbsNote
	for rows.Next() {
		var note pb.AbsNote
		err = rows.Scan(&note.Id, &note.Comment, &note.CreatedAt)
		if err != nil {
			continue
		}
		notes = append(notes, &note)
	}
	return &pb.GetNotesByStudent{Notes: notes}, nil
}
func (r *StudentRepository) CreateNoteForStudent(companyId string, note string, studentId string) (*pb.AbsResponse, error) {
	_, err := r.db.Exec(`INSERT INTO student_note(id , student_id, comment , company_id) values ($1,$2,$3 , $4)`, uuid.New(), studentId, note, companyId)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "created",
	}, nil
}
func (r *StudentRepository) DeleteStudentNote(companyId string, id string) (*pb.AbsResponse, error) {
	_, err := r.db.Exec(`DELETE FROM student_note WHERE id = $1 and company_id=$2`, id, companyId)
	if err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "deleted",
	}, nil
}
func (r *StudentRepository) SearchStudent(companyId string, value string) (*pb.SearchStudentResponse, error) {
	query := `
        SELECT id, name, phone 
        FROM students 
        WHERE company_id=$3 and  name ILIKE $1 OR company_id=$3 and phone ILIKE $2 ;
    `

	rows, err := r.db.Query(query, "%"+value+"%", "%"+value+"%", companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*pb.AbsStudent
	for rows.Next() {
		var student pb.AbsStudent
		if err := rows.Scan(&student.Id, &student.Name, &student.PhoneNumber); err != nil {
			return nil, err
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.SearchStudentResponse{Students: students}, nil
}
func (r *StudentRepository) GetHistoryGroupById(companyId string, groupId string) (*pb.GetHistoryGroupResponse, error) {
	response := &pb.GetHistoryGroupResponse{}
	groupHistoryQuery := `SELECT id, field, old_value, current_value 
                          FROM group_history 
                          WHERE group_id = $1 and company_id=$2 order by created_at desc`

	rows, err := r.db.Query(groupHistoryQuery, groupId, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var history pb.AbsHistory
		if err := rows.Scan(&history.Id, &history.EditedField, &history.OldValue, &history.CurrentValue); err != nil {
			return nil, err
		}
		response.GroupHistory = append(response.GroupHistory, &history)
	}

	studentHistoryQuery := `SELECT s.id, s.name, s.phone, gh.old_condition, gh.current_condition, gh.specific_date, gh.created_at , g.name , g.start_time , g.start_date , g.end_date , g.date_type , gs.condition , coalesce(gh.comment ,'')
                            FROM group_students gs
                            JOIN students s ON gs.student_id = s.id   
                            JOIN group_student_condition_history gh ON gs.id = gh.group_student_id
                            JOIN groups g on g.id=gs.group_id
                            WHERE gs.group_id = $1 and g.company_id=$2
                            order by created_at desc
                            `

	rows, err = r.db.Query(studentHistoryQuery, groupId, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var studentHistory pb.AbsStudentHistory
		var student pb.AbsStudent
		var group pb.AbsGroup
		var createdAt string
		if err := rows.Scan(&student.Id, &student.Name, &student.PhoneNumber, &studentHistory.OldCondition, &studentHistory.CurrentCondition,
			&studentHistory.SpecificDate, &createdAt, &group.Name, &group.LessonStartTime, &group.GroupStartDate, &group.GroupEndDate, &group.DateType, &group.CurrentGroupStatus, &studentHistory.Comment); err != nil {
			return nil, err
		}
		studentHistory.Group = &group
		studentHistory.Student = &student
		studentHistory.CreatedAt = createdAt
		response.StudentsHistory = append(response.StudentsHistory, &studentHistory)
	}

	return response, nil
}
func (r *StudentRepository) GetHistoryByStudentId(companyId string, studentId string) (*pb.GetHistoryStudentResponse, error) {
	response := &pb.GetHistoryStudentResponse{}
	studentHistoryQuery := `SELECT id, field, old_value, current_value 
                            FROM student_history 
                            WHERE student_id = $1 and company_id=$2
                            ORDER BY created_at DESC`

	rows, err := r.db.Query(studentHistoryQuery, studentId, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var history pb.AbsHistory
		if err := rows.Scan(&history.Id, &history.EditedField, &history.OldValue, &history.CurrentValue); err != nil {
			return nil, err
		}
		response.StudentHistory = append(response.StudentHistory, &history)
	}
	conditionsHistoryQuery := `SELECT s.id, s.name, s.phone, gh.old_condition, gh.current_condition, gh.specific_date, gh.created_at, 
                                   g.id, g.name, g.start_time, g.start_date, g.end_date, g.date_type, 
                                   gs.condition, c.id , c.price , c.description , c.title , c.course_duration , g.is_archived , coalesce(gh.comment , '')
                               FROM group_students gs
                               JOIN students s ON gs.student_id = s.id
                               JOIN group_student_condition_history gh ON gs.id = gh.group_student_id
                               JOIN groups g ON g.id = gs.group_id
                               JOIN courses c on c.id=g.course_id
                               WHERE gs.student_id = $1 and gs.company_id=$2
                               ORDER BY gh.created_at DESC`

	rows, err = r.db.Query(conditionsHistoryQuery, studentId, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var studentHistory pb.AbsStudentHistory
		var student pb.AbsStudent
		var group pb.AbsGroup
		var createdAt string
		var course pb.AbsCourse

		if err := rows.Scan(&student.Id, &student.Name, &student.PhoneNumber,
			&studentHistory.OldCondition, &studentHistory.CurrentCondition,
			&studentHistory.SpecificDate, &createdAt,
			&group.Id, &group.Name, &group.LessonStartTime,
			&group.GroupStartDate, &group.GroupEndDate,
			&group.DateType, &group.CurrentGroupStatus, &course.Id, &course.Price, &course.Description, &course.Name, &course.CourseDuration, &group.IsArchived, &studentHistory.Comment); err != nil {
			return nil, err
		}

		group.Course = &course
		studentHistory.Student = &student
		studentHistory.Group = &group
		studentHistory.CreatedAt = createdAt

		response.ConditionsHistory = append(response.ConditionsHistory, &studentHistory)
	}

	return response, nil
}
func (r *StudentRepository) TransferLessonDate(companyId string, groupId string, from string, to string) (*pb.AbsResponse, error) {
	//validDay, err := utils.IsValidLessonDay(r.db, groupId, from)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if !validDay {
	//	return &pb.AbsResponse{
	//		Status:  403,
	//		Message: "The selected 'from' date does not match the group's lesson days",
	//	}, nil
	//}
	var checker bool
	err := r.db.QueryRow(`SELECT exists(SELECT 1 FROM transfer_lesson where group_id=$1 and real_date=$2 and transfer_date=$3 and company_id=$4)`, groupId, from, to, companyId).Scan(&checker)
	if err != nil {
		return nil, err
	}
	if checker {
		_, err = r.db.Exec(`DELETE FROM transfer_lesson where group_id=$1 and real_date=$2 and transfer_date=$3 and company_id=$4`, groupId, from, to, companyId)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = r.db.Exec(`INSERT INTO transfer_lesson(id, group_id, real_date, transfer_date , company_id) values ($1, $2, $3, $4 , $5)`, uuid.New(), groupId, from, to, companyId)
		if err != nil {
			return nil, err
		}
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "accomplished",
	}, nil
}
func (r *StudentRepository) ChangeConditionStudent(ctx context.Context, companyId string, studentId string, groupId string, status string, returnTheMoney bool, tillDate string, actionById, actionByName, comment string) (*pb.AbsResponse, error) {
	isEliminatedInTrial := false

	if err := r.ensureFinanceClient(); err != nil {
		return nil, err
	}

	validStatuses := map[string]bool{"FREEZE": true, "ACTIVE": true, "DELETE": true}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	var tillDateParsed sql.NullTime
	if tillDate != "" {
		parsedDate, err := time.Parse("2006-01-02", tillDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		tillDateParsed = sql.NullTime{Time: parsedDate, Valid: true}
	} else {
		tillDateParsed = sql.NullTime{Valid: false}
	}

	if !tillDateParsed.Valid || tillDateParsed.Time.After(time.Now()) {
		return nil, fmt.Errorf("invalid tillDate: it must be in the past and valid")
	}

	if !r.checkArgumentsIsActive(companyId, groupId, studentId) {
		return nil, fmt.Errorf("group or student is not active")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	var oldCondition, groupStudentId string
	err = tx.QueryRow(`
        SELECT condition, id
        FROM group_students
        WHERE student_id = $1 AND group_id = $2 and company_id=$3`, studentId, groupId, companyId).
		Scan(&oldCondition, &groupStudentId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve old condition: %v", err)
	}

	if oldCondition == status {
		tx.Rollback()
		return nil, fmt.Errorf("condition is the same as you give")
	}

	updateStmt := `
        UPDATE group_students
        SET condition = $1,
            last_specific_date = COALESCE($2, NOW())
        WHERE id=$3 and company_id=$4
    `
	_, err = tx.Exec(updateStmt, status, tillDateParsed, groupStudentId, companyId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update group_students: %v", err)
	}
	if oldCondition == "FREEZE" && status == "DELETE" {
		var exists bool
		err = tx.QueryRow(`
            SELECT EXISTS(
                SELECT 1
                FROM group_student_condition_history
                WHERE student_id = $1 AND group_id = $2 AND group_student_id = $3 and company_id=$3
            )`, studentId, groupId, groupStudentId, companyId).
			Scan(&exists)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to check existence in history: %v", err)
		}
		isEliminatedInTrial = !exists
	}

	insertHistoryStmt := `
        INSERT INTO group_student_condition_history (id, group_student_id, student_id, group_id, old_condition, current_condition, specific_date, return_the_money, created_at, is_eliminated_trial , company_id , comment)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), $9 , $10 , $11)
    `
	_, err = tx.Exec(insertHistoryStmt, uuid.New(), groupStudentId, studentId, groupId, oldCondition, status, tillDate, returnTheMoney, isEliminatedInTrial, companyId, comment)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to insert into group_student_condition_history: %v", err)
	}

	fromDate := tillDateParsed.Time
	currentDate := time.Now()
	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()
	if returnTheMoney {
		for d := fromDate; d.Before(currentDate) || d.Equal(currentDate); d = d.AddDate(0, 1, 0) {
			monthYearDate := d.Format("2006-01-02")

			manaulPriceForCourse, _ := r.financeClient.GetDiscountByStudentId(ctx, studentId, groupId)
			amount, err := utils.CalculateMoneyForStatus(r.db, manaulPriceForCourse, groupId, monthYearDate)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to calculate money for %s: %v", monthYearDate, err)
			}

			var description, transactionType string
			switch status {
			case "FREEZE":
				description = fmt.Sprintf("Student guruhdan muzlatildi. %s oyi uchun pul qaytarib berildi.", monthYearDate)
				transactionType = "REFUND"
			case "DELETE":
				description = fmt.Sprintf("Student guruhdan o'chirildi. %s oyi uchun pul qaytarib berildi.", monthYearDate)
				transactionType = "REFUND"
			case "ACTIVE":
				description = fmt.Sprintf("Student guruhga qo'shildi. %s oyi uchun pul hisoblandi va yechib olindi.", monthYearDate)
				transactionType = "TAKE_OFF"
			default:
				tx.Rollback()
				return nil, fmt.Errorf("unknown status: %s", status)
			}

			_, err = r.financeClient.PaymentAdd(ctx,
				description, monthYearDate, "CASH", fmt.Sprintf("%v", amount),
				studentId, transactionType, actionById, actionByName, groupId, tillDate)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to add payment for %s: %v", monthYearDate, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &pb.AbsResponse{
		Message: "Condition changed successfully",
		Status:  200,
	}, nil
}
func (r *StudentRepository) GetStudentsByGroupId(companyId string, groupId string, withOutdated bool) (*pb.GetStudentsByGroupIdResponse, error) {
	var students []*pb.AbsStudent
	query := `
        SELECT s.id, s.name, s.phone
        FROM students s
        JOIN group_students gs ON s.id = gs.student_id
        WHERE gs.group_id = $1 and s.company_id=$2
    `

	if !withOutdated {
		query += " AND gs.condition = 'ACTIVE'"
	}

	rows, err := r.db.Query(query, groupId, companyId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var student pb.AbsStudent
		if err := rows.Scan(&student.Id, &student.Name, &student.PhoneNumber); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return &pb.GetStudentsByGroupIdResponse{Students: students}, nil
}
func (r *StudentRepository) ChangeUserBalanceHistory(ctx context.Context, companyId string, comment string, groupId string, createdById string, createdByName string, givenDate string, amount string, paymentType string, studentId string) (*pb.AbsResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to begin transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	var currentBalance float64
	var groupName string

	err = tx.QueryRow("SELECT balance FROM students WHERE id = $1 and company_id=$2", studentId, companyId).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Failed to get current balance: %v", err)
	}

	if paymentType != "TAKE_OFF" && groupId != "" {
		err = tx.QueryRow("SELECT name FROM groups WHERE id = $1 and company_id=$2", groupId, companyId).Scan(&groupName)
		if errors.Is(err, sql.ErrNoRows) {
			groupName = ""
		} else if err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "Failed to query group name: %v", err)
		}
	}

	amountValue, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "Invalid amount: %v", err)
	}

	var newBalance float64
	switch paymentType {
	case "ADD":
		newBalance = currentBalance + amountValue
	case "TAKE_OFF":
		newBalance = currentBalance - amountValue
	default:
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment type: %s", paymentType)
	}
	err = r.BalanceHistoryMaker(companyId, tx, currentBalance, newBalance, studentId, comment, groupId, groupName, createdById, createdByName, givenDate, amount, paymentType)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	if paymentType == "ADD" {
		r.SendSmsPaymentSuccessful(ctx, studentId, groupId, companyId, amountValue)
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "balance edited",
	}, nil
}
func (r *StudentRepository) ChangeUserBalanceHistoryByDebit(companyId string, studentId string, oldDebit string, currentDebit string, givenDate string, comment string, paymentType string, createdById string, createdByName string, groupId string) (*pb.AbsResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to begin transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	var currentBalance float64
	var groupName string

	err = tx.QueryRow("SELECT balance FROM students WHERE id = $1", studentId).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Failed to get current balance: %v", err)
	}

	if paymentType != "TAKE_OFF" && groupId != "" {
		err = tx.QueryRow("SELECT name FROM groups WHERE id = $1", groupId).Scan(&groupName)
		if errors.Is(err, sql.ErrNoRows) {
			groupName = ""
		} else if err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "Failed to query group name: %v", err)
		}
	}
	oldBalance := currentBalance
	amountValue, err := strconv.ParseFloat(oldDebit, 64)
	currentAmountValue, err := strconv.ParseFloat(currentDebit, 64)
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "Invalid amount: %v", err)
	}
	switch paymentType {
	case "ADD":
		currentBalance = (currentBalance - amountValue) + currentAmountValue
	case "TAKE_OFF":
		currentBalance = (currentBalance + amountValue) - currentAmountValue
	default:
		tx.Rollback()
		return nil, status.Errorf(codes.Aborted, "invalid payment type")
	}
	err = r.BalanceHistoryMaker(companyId, tx, oldBalance, currentBalance, studentId, comment, groupId, groupName, createdById, createdByName, givenDate, fmt.Sprintf("%.2f", currentAmountValue), paymentType)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return &pb.AbsResponse{
		Status:  http.StatusOK,
		Message: "balance edited",
	}, nil
}
func (r *StudentRepository) BalanceHistoryMaker(companyId string, tx *sql.Tx, currentBalance, newBalance float64, studentId string, comment, groupId, groupName, createdById, createdByName, givenDate, amount, paymentType string) error {
	result, err := tx.Exec("UPDATE students SET balance = $1 WHERE id = $2", newBalance, studentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return err
	}
	if groupId != "" && groupName == "" {
		groupName = "Group"
	}
	historyData := map[string]interface{}{
		"comment":       comment,
		"groupId":       groupId,
		"groupName":     groupName,
		"createdById":   createdById,
		"createdByName": createdByName,
		"givenDate":     givenDate,
		"amount":        amount,
		"paymentType":   paymentType,
	}

	historyJSON, err := json.Marshal(historyData)
	if err != nil {
		tx.Rollback()
		return err
	}

	field := "balance_add"
	if paymentType == "TAKE_OFF" {
		field = "balance_take_off"
	}
	_, err = tx.Exec(`
        INSERT INTO student_history (id, student_id, field, old_value, current_value, created_at , company_id)
        VALUES (gen_random_uuid(), $1, $2, $3, $4, $5 , $6)
    `, studentId, field, currentBalance, historyJSON, time.Now(), companyId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
func (r *StudentRepository) checkArgumentsIsActive(companyId string, groupId, studentId string) bool {
	var checker bool
	err := r.db.QueryRow(`SELECT exists(SELECT 1 FROM groups where is_archived='false' and id=$1 and company_id=$2)`, groupId, companyId).Scan(&checker)
	if err != nil || !checker {
		return false
	}
	err = r.db.QueryRow(`SELECT exists(SELECT 1 FROM students where condition='ACTIVE' and id=$1 and company_id=$2)`, studentId, companyId).Scan(&checker)
	if err != nil || !checker {
		return false
	}

	return checker
}
func (r *StudentRepository) StudentBalanceTaker() {
	rows, err := r.db.Query(`SELECT id FROM company WHERE valid_date > CURRENT_DATE`)
	if err != nil {
		fmt.Println("error while getting companies:", err)
		return
	}

	for rows.Next() {
		var companyId string
		err := rows.Scan(&companyId)
		if err != nil {
			fmt.Println("erorr while scanning companyid")
			continue
		}
		if err := r.ensureFinanceClient(); err != nil {
			return
		}
		rows, err := r.db.Query(`SELECT id FROM students where condition='ACTIVE' and company_id = $1 `, companyId)
		if err != nil {
			fmt.Printf("error get active student %v", err)
			return
		}
		defer rows.Close()
		ctx, cancelFunc := utils.NewTimoutContext(context.Background(), companyId)
		defer cancelFunc()
		for rows.Next() {
			var studentId string
			err = rows.Scan(&studentId)
			if err != nil {
				fmt.Printf("error scanning active student %v", err)
				continue
			}
			extraRow, err := r.db.Query(`SELECT group_id FROM group_students where student_id=$1 and condition='ACTIVE' and company_id= $2`, studentId, companyId)
			fmt.Println(studentId)
			if err != nil {
				fmt.Printf("error getting  active groupid %v", err)
				continue
			}
			for extraRow.Next() {
				var (
					groupId     string
					takingPrice float64
					comment     string
				)
				err = extraRow.Scan(&groupId)
				if err != nil {
					fmt.Printf("error scanning groupid student %v", err)
					continue
				}
				discountAmount, _ := r.financeClient.GetDiscountByStudentId(ctx, studentId, groupId)
				err = r.db.QueryRow(`SELECT c.price FROM groups g join courses c on g.course_id=c.id where g.id=$1 and c.company_id=$2`, groupId, companyId).Scan(&takingPrice)
				if err != nil {
					fmt.Printf("error getting course price active student %v", err)
					continue
				}
				if discountAmount == nil {
					comment = "ushbu oy uchun oylik tolov student balansidan yechib olindi."
				} else {
					takingPrice = takingPrice - *discountAmount
					comment = "ushbu oy uchun oylik tolov student balansidan yechib olindi chegirma narxida"
				}
				//_, err := r.ChangeUserBalanceHistory("ushbu oy uchun oylik tolov student balansidan yechib olindi.", groupId, "00000000-0000-0000-0000-000000000000", "TIZIM", time.Now().Format("2006-01-02"), takingPrice, "TAKE_OFF", studentId)
				_, err :=
					r.financeClient.PaymentAdd(ctx, comment, time.Now().Format("2006-01-02"), "CASH", fmt.Sprintf("%.2f", takingPrice), studentId, "TAKE_OFF", "00000000-0000-0000-0000-000000000000", "TIZIM", groupId, time.Now().AddDate(0, 0, -1).String())
				if err != nil {
					fmt.Printf("error changing balance history active student %v", err)
					continue
				}
			}
			extraRow.Close()
		}
	}
}

func (r *StudentRepository) CalculateDiscountSumma(companyId string, groupId string, startDate string, endDate string, discountPrice string, studentId string, paymentDate, activationDate string) (*pb.CalculateDiscountResponse, error) {
	groupIDInt, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %v", err)
	}

	discountPriceFloat, err := strconv.ParseFloat(discountPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid discount price: %v", err)
	}

	var group struct {
		EndDate string
		Days    []string
	}

	err = r.db.QueryRow(`
        SELECT end_date, days 
        FROM groups 
        WHERE id = $1 AND company_id = $2
    `, groupIDInt, companyId).Scan(&group.EndDate, pq.Array(&group.Days))

	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %v", err)
	}
	paymentDateTime, err := parseDate(paymentDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment date: %v", err)
	}

	activationDateTime, err := parseDate(activationDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse activation date: %v", err)
	}

	groupEndDateTime, err := parseDate(group.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse group end date: %v", err)
	}

	monthStart := time.Date(paymentDateTime.Year(), paymentDateTime.Month(), 1, 0, 0, 0, 0, paymentDateTime.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)

	totalMonthLessons, err := r.countLessonsBetweenDates(
		monthStart.Format("2006-01-02"),
		monthEnd.Format("2006-01-02"),
		group.Days,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count total month lessons: %v", err)
	}

	remainingStart := monthStart
	remainingEnd := monthEnd

	if groupEndDateTime.Before(monthEnd) && groupEndDateTime.After(monthStart) {
		remainingEnd = groupEndDateTime
	}

	if activationDateTime.After(monthStart) && activationDateTime.Before(monthEnd) {
		remainingStart = activationDateTime
	}

	remainingLessons, err := r.countLessonsBetweenDates(
		remainingStart.Format("2006-01-02"),
		remainingEnd.Format("2006-01-02"),
		group.Days,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count remaining lessons: %v", err)
	}

	finalDiscountAmount := discountPriceFloat * float64(remainingLessons) / float64(totalMonthLessons)

	return &pb.CalculateDiscountResponse{
		CalculatedPrice: cast.ToString(finalDiscountAmount),
	}, nil
}
func parseDate(input string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02T15:04:05Z", input); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", input)
}

func (r *StudentRepository) countLessonsBetweenDates(startDate, endDate string, groupDays []string) (int, error) {
	start, err := parseDate(startDate)
	if err != nil {
		return 0, err
	}

	end, err := parseDate(endDate)
	if err != nil {
		return 0, err
	}

	weekdays := make(map[string]time.Weekday)
	weekdays["DUSHANBA"] = time.Monday
	weekdays["SESHANBA"] = time.Tuesday
	weekdays["CHORSHANBA"] = time.Wednesday
	weekdays["PAYSHANBA"] = time.Thursday
	weekdays["JUMA"] = time.Friday
	weekdays["SHANBA"] = time.Saturday
	weekdays["YAKSHANBA"] = time.Sunday

	groupWeekdays := make(map[time.Weekday]bool)
	for _, day := range groupDays {
		if wd, ok := weekdays[day]; ok {
			groupWeekdays[wd] = true
		}
	}

	lessonCount := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if groupWeekdays[d.Weekday()] {
			lessonCount++
		}
	}

	return lessonCount, nil
}

func (r *StudentRepository) SendSmsReasonAddToGroup(ctx context.Context, studentId string, groupId string, companyId string) {
	var smsTemplate models.SmsTemplate
	var texts []string
	var teacherId, phoneNumber string
	var smsBalance int
	teacherName := "O'qituvchi"

	// Fetch SMS template
	err := r.db.QueryRow(`SELECT id, texts, action_type, sms_template_type, is_active FROM sms_template WHERE company_id=$1 AND action_type='JOINED_GROUP_ALERT' AND sms_template_type='ACTION' AND is_active=true`, companyId).
		Scan(&smsTemplate.ID, pq.Array(&texts), &smsTemplate.ActionType, &smsTemplate.SmsTemplateType, &smsTemplate.IsActive)
	if err != nil {
		fmt.Println("error fetching sms template:", err)
		return
	}
	fmt.Println("Fetched SMS template successfully")

	// Fetch initial SMS balance
	err = r.db.QueryRow(`SELECT sms_balance FROM company WHERE id=$1`, companyId).Scan(&smsBalance)
	if err != nil {
		fmt.Println("error fetching SMS balance:", err)
		return
	}
	fmt.Printf("Initial SMS balance: %d\n", smsBalance)
	if smsBalance < 1 {
		fmt.Println("not enough SMS balance")
		return
	}

	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()

	// Get teacher and phone number
	if err = r.db.QueryRow(`SELECT teacher_id FROM groups where id=$1`, groupId).Scan(&teacherId); err == nil {
		teacherName, _ = r.userClient.GetTeacherById(ctx, teacherId)
	}
	err = r.db.QueryRow(`SELECT phone FROM students where id=$1`, studentId).Scan(&phoneNumber)
	if err != nil {
		fmt.Println("error fetching phone number:", err)
		return
	}
	fmt.Println("Teacher and phone fetched successfully")

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Println("error starting transaction:", err)
		return
	}
	fmt.Println("Transaction started")

	// Re-check SMS balance inside transaction
	err = tx.QueryRow(`SELECT sms_balance FROM company WHERE id=$1`, companyId).Scan(&smsBalance)
	if err != nil {
		fmt.Println("error fetching SMS balance (transaction):", err)
		tx.Rollback()
		return
	}
	fmt.Printf("SMS balance inside transaction: %d\n", smsBalance)
	if smsBalance < 1 {
		fmt.Println("not enough SMS balance (transaction)")
		tx.Rollback()
		return
	}

	// Format SMS text
	smsText, usedSmsCount := utils.GetSmsFormatted(strings.Join(texts, " "), teacherName, r.db, studentId, groupId, 0)
	fmt.Printf("SMS formatted successfully: %s\n", smsText)

	// Insert into sms_used
	_, err = tx.Exec(`INSERT INTO sms_used(id, company_id, sms_template_id, texts, sms_count, student_id, sms_used_type, created_at, created_by_id, created_by_name) 
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		uuid.New(), companyId, smsTemplate.ID, pq.Array(texts), usedSmsCount, studentId, "BY_TEMPLATE", time.Now(), "00000000-0000-0000-0000-000000000000", "system")
	if err != nil {
		fmt.Println("error inserting into sms_used:", err)
		tx.Rollback()
		return
	}
	fmt.Println("Inserted into sms_used successfully")

	// Update company balance
	_, err = tx.Exec(`UPDATE company SET sms_balance=sms_balance-$1 WHERE id=$2`, usedSmsCount, companyId)
	if err != nil {
		fmt.Println("error updating company balance:", err)
		tx.Rollback()
		return
	}
	fmt.Println("Company balance updated successfully")

	// Send SMS after transaction but before commit
	err = utils.SendSMS(phoneNumber, smsText)
	if err != nil {
		fmt.Println("error sending SMS after transaction:", err)
		tx.Rollback()
		return
	}
	fmt.Println("SMS sent successfully")

	// Commit transaction
	if err := tx.Commit(); err != nil {
		fmt.Println("error committing transaction:", err)
		return
	}
	fmt.Println("Transaction committed")
}

func (r *StudentRepository) SendSmsPaymentSuccessful(ctx context.Context, studentId string, groupId string, companyId string, amountValue float64) {
	var smsTemplate models.SmsTemplate
	var texts []string
	var teacherId, phoneNumber string
	var smsBalance int
	teacherName := "O'qituvchi"

	// Fetch SMS template
	err := r.db.QueryRow(`SELECT id, texts, action_type, sms_template_type, is_active FROM sms_template WHERE company_id=$1 AND action_type='PAYMENT_SUCCESSFUL_ALERT' AND sms_template_type='ACTION' AND is_active=true`, companyId).
		Scan(&smsTemplate.ID, pq.Array(&texts), &smsTemplate.ActionType, &smsTemplate.SmsTemplateType, &smsTemplate.IsActive)
	if err != nil {
		fmt.Println("error fetching sms template:", err)
		return
	}
	fmt.Println("Fetched SMS template successfully")

	// Fetch initial SMS balance
	err = r.db.QueryRow(`SELECT sms_balance FROM company WHERE id=$1`, companyId).Scan(&smsBalance)
	if err != nil {
		fmt.Println("error fetching SMS balance:", err)
		return
	}
	fmt.Printf("Initial SMS balance: %d\n", smsBalance)
	if smsBalance < 1 {
		fmt.Println("not enough SMS balance")
		return
	}

	ctx, cancelFunc := utils.NewTimoutContext(ctx, companyId)
	defer cancelFunc()

	// Get teacher and phone number
	if err = r.db.QueryRow(`SELECT teacher_id FROM groups where id=$1`, groupId).Scan(&teacherId); err == nil {
		teacherName, _ = r.userClient.GetTeacherById(ctx, teacherId)
	}
	err = r.db.QueryRow(`SELECT phone FROM students where id=$1`, studentId).Scan(&phoneNumber)
	if err != nil {
		fmt.Println("error fetching phone number:", err)
		return
	}
	fmt.Println("Teacher and phone fetched successfully")

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Println("error starting transaction:", err)
		return
	}
	fmt.Println("Transaction started")

	// Re-check SMS balance inside transaction
	err = tx.QueryRow(`SELECT sms_balance FROM company WHERE id=$1`, companyId).Scan(&smsBalance)
	if err != nil {
		fmt.Println("error fetching SMS balance (transaction):", err)
		tx.Rollback()
		return
	}
	fmt.Printf("SMS balance inside transaction: %d\n", smsBalance)
	if smsBalance < 1 {
		fmt.Println("not enough SMS balance (transaction)")
		tx.Rollback()
		return
	}

	// Format SMS text
	smsText, usedSmsCount := utils.GetSmsFormatted(strings.Join(texts, " "), teacherName, r.db, studentId, groupId, amountValue)
	fmt.Printf("SMS formatted successfully: %s\n", smsText)

	// Insert into sms_used
	_, err = tx.Exec(`INSERT INTO sms_used(id, company_id, sms_template_id, texts, sms_count, student_id, sms_used_type, created_at, created_by_id, created_by_name) 
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		uuid.New(), companyId, smsTemplate.ID, pq.Array(texts), usedSmsCount, studentId, "BY_TEMPLATE", time.Now(), "00000000-0000-0000-0000-000000000000", "system")
	if err != nil {
		fmt.Println("error inserting into sms_used:", err)
		tx.Rollback()
		return
	}
	fmt.Println("Inserted into sms_used successfully")

	// Update company balance
	_, err = tx.Exec(`UPDATE company SET sms_balance=sms_balance-$1 WHERE id=$2`, usedSmsCount, companyId)
	if err != nil {
		fmt.Println("error updating company balance:", err)
		tx.Rollback()
		return
	}
	fmt.Println("Company balance updated successfully")

	// Send SMS after transaction but before commit
	err = utils.SendSMS(phoneNumber, smsText)
	if err != nil {
		fmt.Println("error sending SMS after transaction:", err)
		tx.Rollback()
		return
	}
	fmt.Println("SMS sent successfully")

	// Commit transaction
	if err := tx.Commit(); err != nil {
		fmt.Println("error committing transaction:", err)
		return
	}
	fmt.Println("Transaction committed")
}
