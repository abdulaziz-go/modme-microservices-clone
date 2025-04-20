package client

import (
	"api-gateway/grpc/proto/pb"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

type EducationClient struct {
	roomClient           pb.RoomServiceClient
	courseClient         pb.CourseServiceClient
	groupClient          pb.GroupServiceClient
	attendanceClient     pb.AttendanceServiceClient
	studentClient        pb.StudentServiceClient
	companyClient        pb.CompanyServiceClient
	tariffClient         pb.TariffServiceClient
	companyFinanceClient pb.CompanyFinanceServiceClient
	smsServiceClient     pb.SmsServiceClient
}

func NewEducationClient(addr string) (*EducationClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	roomClient := pb.NewRoomServiceClient(conn)
	courseClient := pb.NewCourseServiceClient(conn)
	groupClient := pb.NewGroupServiceClient(conn)
	attendanceClient := pb.NewAttendanceServiceClient(conn)
	studentClient := pb.NewStudentServiceClient(conn)
	companyClient := pb.NewCompanyServiceClient(conn)
	tariffClient := pb.NewTariffServiceClient(conn)
	companyFinanceClient := pb.NewCompanyFinanceServiceClient(conn)
	smsServiceClient := pb.NewSmsServiceClient(conn)
	return &EducationClient{roomClient: roomClient, courseClient: courseClient, groupClient: groupClient, attendanceClient: attendanceClient, studentClient: studentClient, companyClient: companyClient, tariffClient: tariffClient, companyFinanceClient: companyFinanceClient, smsServiceClient: smsServiceClient}, nil
}

// Education Service method client
func (lc *EducationClient) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.AbsResponse, error) {
	return lc.roomClient.CreateRoom(ctx, req)
}

func (lc *EducationClient) UpdateRoom(ctx context.Context, req *pb.AbsRoom) (*pb.AbsResponse, error) {
	return lc.roomClient.UpdateRoom(ctx, req)
}

func (lc *EducationClient) DeleteRoom(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := pb.DeleteAbsRequest{
		Id: id,
	}
	return lc.roomClient.DeleteRoom(ctx, &req)
}

func (lc *EducationClient) GetRoom(ctx context.Context) (*pb.GetUpdateRoomAbs, error) {
	return lc.roomClient.GetRooms(ctx, &emptypb.Empty{})
}

func (lc *EducationClient) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.AbsResponse, error) {
	return lc.courseClient.CreateCourse(ctx, req)
}

func (lc *EducationClient) UpdateCourse(ctx context.Context, req *pb.AbsCourse) (*pb.AbsResponse, error) {
	return lc.courseClient.UpdateCourse(ctx, req)
}

func (lc *EducationClient) DeleteCourse(ctx context.Context, id string) (*pb.AbsResponse, error) {
	req := pb.DeleteAbsRequest{
		Id: id,
	}
	return lc.courseClient.DeleteCourse(ctx, &req)
}

func (lc *EducationClient) GetCourse(ctx context.Context) (*pb.GetUpdateCourseAbs, error) {
	return lc.courseClient.GetCourses(ctx, &emptypb.Empty{})
}

func (lc *EducationClient) GetCourseById(ctx context.Context, id string) (*pb.GetCourseByIdResponse, error) {
	return lc.courseClient.GetCourseById(ctx, &pb.GetCourseByIdRequest{Id: id})
}

func (lc *EducationClient) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.AbsResponse, error) {
	return lc.groupClient.CreateGroup(ctx, req)
}

func (lc *EducationClient) UpdateGroup(ctx context.Context, req *pb.GetUpdateGroupAbs) (*pb.AbsResponse, error) {
	return lc.groupClient.UpdateGroup(ctx, req)
}

func (lc *EducationClient) DeleteGroup(ctx context.Context, id string) (*pb.AbsResponse, error) {
	return lc.groupClient.DeleteGroup(ctx, &pb.DeleteAbsRequest{Id: id})
}

func (lc *EducationClient) GetAllGroup(ctx context.Context, req *pb.GetGroupsRequest) (*pb.GetGroupsResponse, error) {
	return lc.groupClient.GetGroups(ctx, req)
}

func (lc *EducationClient) GetGroupById(ctx context.Context, groupId, userId, role string) (*pb.GetGroupAbsResponse, error) {
	return lc.groupClient.GetGroupById(ctx, &pb.GetGroupByIdRequest{Id: groupId, ActionId: userId, ActionRole: role})
}

func (lc *EducationClient) GetAttendanceByGroup(ctx context.Context, req *pb.GetAttendanceRequest) (*pb.GetAttendanceResponse, error) {
	return lc.attendanceClient.GetAttendance(ctx, req)
}

func (lc *EducationClient) SetAttendanceByGroup(ctx context.Context, req *pb.SetAttendanceRequest) (*pb.AbsResponse, error) {
	return lc.attendanceClient.SetAttendance(ctx, req)
}

func (lc *EducationClient) GetGroupByCourseId(ctx context.Context, courseId string) (*pb.GetGroupsByCourseResponse, error) {
	return lc.groupClient.GetGroupsByCourseId(ctx, &pb.GetGroupByIdRequest{Id: courseId})
}

func (lc *EducationClient) GetAllStudent(ctx context.Context, req *pb.GetAllStudentRequest) (*pb.GetAllStudentResponse, error) {
	response, err := lc.studentClient.GetAllStudent(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (lc *EducationClient) CreateStudent(ctx context.Context, p *pb.CreateStudentRequest) (*pb.AbsResponse, error) {
	response, err := lc.studentClient.CreateStudent(ctx, p)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (lc *EducationClient) AddStudentToGroup(ctx context.Context, p *pb.AddToGroupRequest) (*pb.AbsResponse, error) {
	return lc.studentClient.AddToGroup(ctx, p)
}

func (lc *EducationClient) UpdateStudent(ctx context.Context, p *pb.UpdateStudentRequest) (*pb.AbsResponse, error) {
	return lc.studentClient.UpdateStudent(ctx, p)
}

func (lc *EducationClient) DeleteStudent(ctx context.Context, id string, returnMoney bool, actionId, actionName string) (*pb.AbsResponse, error) {
	return lc.studentClient.DeleteStudent(ctx, &pb.DeleteStudentRequest{StudentId: id, ReturnMoney: returnMoney, ActionByName: actionName, ActionById: actionId})
}

func (lc *EducationClient) GetStudentById(ctx context.Context, id string) (*pb.GetStudentByIdResponse, error) {
	return lc.studentClient.GetStudentById(ctx, &pb.NoteStudentByAbsRequest{Id: id})
}

func (lc *EducationClient) GetNotesByStudentId(ctx context.Context, id string) (*pb.GetNotesByStudent, error) {
	return lc.studentClient.GetNoteByStudent(ctx, &pb.NoteStudentByAbsRequest{Id: id})
}

func (lc *EducationClient) CreateNoteForStudent(ctx context.Context, p *pb.CreateNoteRequest) (*pb.AbsResponse, error) {
	return lc.studentClient.CreateNoteForStudent(ctx, p)
}

func (lc *EducationClient) DeleteNote(ctx context.Context, note string) (*pb.AbsResponse, error) {
	return lc.studentClient.DeleteStudentNote(ctx, &pb.NoteStudentByAbsRequest{Id: note})
}

func (lc *EducationClient) SearchStudentByPhoneName(ctx context.Context, value string) (*pb.SearchStudentResponse, error) {
	return lc.studentClient.SearchStudent(ctx, &pb.SearchStudentRequest{Value: value})
}

func (lc *EducationClient) GetHistoryGroupById(ctx context.Context, value string) (*pb.GetHistoryGroupResponse, error) {
	return lc.studentClient.GetHistoryGroupById(ctx, &pb.NoteStudentByAbsRequest{Id: value})
}

func (lc *EducationClient) GetHistoryStudentById(ctx context.Context, value string) (*pb.GetHistoryStudentResponse, error) {
	return lc.studentClient.GetHistoryStudentById(ctx, &pb.NoteStudentByAbsRequest{Id: value})
}

func (lc *EducationClient) TransferLessonDate(ctx context.Context, p *pb.TransferLessonRequest) (*pb.AbsResponse, error) {
	return lc.studentClient.TransferLessonDate(ctx, p)
}

func (lc *EducationClient) ChangeConditionStudent(ctx context.Context, p *pb.ChangeConditionStudentRequest) (*pb.AbsResponse, error) {
	return lc.studentClient.ChangeConditionStudent(ctx, p)
}

func (lc *EducationClient) GetInformationByTeacher(ctx context.Context, teacherId string, isArchived bool) (*pb.GetGroupsByTeacherResponse, error) {
	return lc.groupClient.GetGroupsByTeacherId(ctx, &pb.GetGroupsByTeacherIdRequest{
		TeacherId:  teacherId,
		IsArchived: isArchived,
	})
}

func (lc *EducationClient) GetCommonEducationInformation(ctx context.Context) (int, int, int, int, int) {
	response, err := lc.groupClient.GetCommonInformationEducation(ctx, &emptypb.Empty{})

	if err != nil {
		return 0, 0, 0, 0, 0
	}
	return int(response.ActiveStudentCount), int(response.ActiveGroupCount), int(response.LeaveGroupCount), int(response.DebtorsCount), int(response.EleminatedInTrial)
}

func (lc *EducationClient) CalculateSalaryByTeacher(ctx context.Context, from string, to string, teacherId string) (*pb.CalculateTeacherSalaryResponse, error) {
	return lc.attendanceClient.CalculateTeacherSalaryByAttendance(ctx, &pb.CalculateTeacherSalaryRequest{
		From:      from,
		To:        to,
		TeacherId: teacherId,
	})
}

func (lc *EducationClient) GetLeftAfterTrialPeriod(ctx context.Context, from string, to string, page string, size string) (*pb.GetLeftAfterTrialPeriodResponse, error) {
	return lc.groupClient.GetLeftAfterTrialPeriod(ctx, &pb.GetLeftAfterTrialPeriodRequest{
		From: from,
		To:   to,
		Page: page,
		Size: size,
	})
}

func (lc *EducationClient) GetCompanyBySubdomain(domain string) (*pb.GetCompanyResponse, error) {
	return lc.companyClient.GetCompanyBySubdomain(context.TODO(), &pb.GetCompanyRequest{Domain: domain})
}

func (lc *EducationClient) CreateCompanyRequest(req *pb.CreateCompanyRequest) (*pb.AbsResponse, error) {
	return lc.companyClient.CreateCompany(context.TODO(), req)
}

func (lc *EducationClient) GetAllCompanies(page string, size string, filter string) (*pb.GetAllResponse, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sizeInt = 20
	}
	return lc.companyClient.GetAll(context.TODO(), &pb.PageRequest{
		Page:   int32(pageInt),
		Size:   int32(sizeInt),
		Filter: filter,
	})
}

func (lc *EducationClient) UpdateCompany(p *pb.UpdateCompanyRequest) (*pb.AbsResponse, error) {
	return lc.companyClient.UpdateCompany(context.TODO(), p)
}

func (lc *EducationClient) CreateTariff(req *pb.Tariff) (*pb.Tariff, error) {
	return lc.tariffClient.Create(context.TODO(), req)
}

func (lc *EducationClient) UpdateTariff(req *pb.Tariff) (*pb.Tariff, error) {
	return lc.tariffClient.Update(context.TODO(), req)
}

func (lc *EducationClient) DeleteTariff(id int32) (*pb.Tariff, error) {
	return lc.tariffClient.Delete(context.TODO(), &pb.Tariff{Id: id})
}

func (lc *EducationClient) GetAllTariff() *pb.TariffList {
	resp, err := lc.tariffClient.Get(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return nil
	}
	return resp
}

func (lc *EducationClient) FinanceCreate(req *pb.CompanyFinance) (*pb.CompanyFinance, error) {
	return lc.companyFinanceClient.Create(context.TODO(), req)
}

func (lc *EducationClient) FinanceDelete(id string) (*pb.AbsResponse, error) {
	return lc.companyFinanceClient.Delete(context.TODO(), &pb.DeleteAbsRequest{Id: id})
}

func (lc *EducationClient) FinanceGetByCompany(req *pb.PageRequest) (*pb.CompanyFinanceSelfList, error) {
	return lc.companyFinanceClient.GetByCompany(context.TODO(), req)
}

func (lc *EducationClient) FinanceGetAll(req *pb.PageRequest) (*pb.CompanyFinanceList, error) {
	return lc.companyFinanceClient.GetAll(context.TODO(), req)
}

func (lc *EducationClient) FinanceUpdate(req *pb.CompanyFinance) (*pb.CompanyFinance, error) {
	return lc.companyFinanceClient.UpdateByCompany(context.TODO(), req)
}

func (lc *EducationClient) GetStatisticCompany(req *pb.GetStatisticRequest) (*pb.GetStatisticResponse, error) {
	return lc.companyClient.GetStatistic(context.TODO(), req)
}

func (lc *EducationClient) GetSmsLogs(ctx context.Context, req *pb.GetSmsLogRequest) (*pb.GetSmsLogResponse, error) {
	return lc.smsServiceClient.GetSmsLogs(ctx, req)
}
func (lc *EducationClient) AddSms(ctx context.Context, req *pb.AddSmsRequest) (*pb.AbsResponse, error) {
	return lc.smsServiceClient.AddSms(ctx, req)
}
func (lc *EducationClient) DeleteSms(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	return lc.smsServiceClient.DeleteSms(ctx, req)
}
func (lc *EducationClient) GetSmsTransactionDetail(ctx context.Context, req *pb.PageRequest) (*pb.GetSmsTransactionDetailResponse, error) {
	return lc.smsServiceClient.GetSmsTransactionDetail(ctx, req)
}
func (lc *EducationClient) GetSmsTemplate(ctx context.Context, req *pb.GetSmsTemplateRequest) (*pb.GetSmsTemplateResponse, error) {
	return lc.smsServiceClient.GetSmsTemplate(ctx, req)
}
func (lc *EducationClient) SetSmsTemplate(ctx context.Context, req *pb.SetSmsTemplateRequest) (*pb.AbsResponse, error) {
	return lc.smsServiceClient.SetSmsTemplate(ctx, req)
}
func (lc *EducationClient) SendSmsDirectly(ctx context.Context, req *pb.SendSmsDirectlyRequest) (*pb.AbsResponse, error) {
	return lc.smsServiceClient.SendSmsDirectly(ctx, req)
}
