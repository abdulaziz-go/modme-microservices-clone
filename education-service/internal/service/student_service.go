package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StudentService struct {
	pb.UnimplementedStudentServiceServer
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}
func (s *StudentService) GetAllStudent(ctx context.Context, req *pb.GetAllStudentRequest) (*pb.GetAllStudentResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetAllStudent(ctx, companyId, req.Condition, req.Page, req.Size, req.CourseId, req.NameOrSurname, req.TrainingDatesStart, req.TrainingDatesEnd)
}
func (s *StudentService) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := s.repo.CreateStudent(companyId, req.CreatedBy, req.PhoneNumber, req.Name, req.GroupId, req.Address, req.AdditionalContact, req.DateFrom, req.DateOfBirth, req.Gender, req.PassportId, req.TelegramUsername); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "student created successfully",
	}, nil
}
func (s *StudentService) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := s.repo.UpdateStudent(companyId, req.StudentId, req.PhoneNumber, req.Name, req.Address, req.AdditionalContact, req.DateOfBirth, req.Gender, req.PassportId); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "student updated successfully",
	}, nil
}
func (s *StudentService) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := s.repo.DeleteStudent(ctx, companyId, req.StudentId, req.ReturnMoney, req.ActionById, req.ActionByName); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "accomplished",
	}, nil
}
func (s *StudentService) AddToGroup(ctx context.Context, req *pb.AddToGroupRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	if err := s.repo.AddToGroup(ctx, companyId, req.GroupId, req.StudentIds, req.CreatedDate, req.CreatedBy); err != nil {
		return nil, err
	}
	return &pb.AbsResponse{
		Status:  200,
		Message: "students added to group",
	}, nil
}
func (s *StudentService) GetStudentById(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetStudentByIdResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetStudentById(ctx, companyId, req.Id)
}
func (s *StudentService) GetNoteByStudent(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetNotesByStudent, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetNoteByStudent(companyId, req.Id)
}
func (s *StudentService) CreateNoteForStudent(ctx context.Context, req *pb.CreateNoteRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.CreateNoteForStudent(companyId, req.Note, req.StudentId)
}
func (s *StudentService) DeleteStudentNote(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.DeleteStudentNote(companyId, req.Id)
}
func (s *StudentService) SearchStudent(ctx context.Context, req *pb.SearchStudentRequest) (*pb.SearchStudentResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.SearchStudent(companyId, req.Value)
}
func (s *StudentService) GetHistoryGroupById(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetHistoryGroupResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetHistoryGroupById(companyId, req.Id)
}
func (s *StudentService) GetHistoryStudentById(ctx context.Context, req *pb.NoteStudentByAbsRequest) (*pb.GetHistoryStudentResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetHistoryByStudentId(companyId, req.Id)
}
func (s *StudentService) TransferLessonDate(ctx context.Context, req *pb.TransferLessonRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.TransferLessonDate(companyId, req.GroupId, req.From, req.To)
}
func (s *StudentService) ChangeConditionStudent(ctx context.Context, req *pb.ChangeConditionStudentRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.ChangeConditionStudent(ctx, companyId, req.StudentId, req.GroupId, req.Status, req.ReturnTheMoney, req.TillDate, req.ActionById, req.ActionByName, req.Comment)
}
func (s *StudentService) GetStudentsByGroupId(ctx context.Context, req *pb.GetStudentsByGroupIdRequest) (*pb.GetStudentsByGroupIdResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.GetStudentsByGroupId(companyId, req.GroupId, req.WithOutdated)
}
func (s *StudentService) ChangeUserBalanceHistory(ctx context.Context, req *pb.ChangeUserBalanceHistoryRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.ChangeUserBalanceHistory(companyId, req.Comment, req.GroupId, req.CreatedBy, req.CreatedByName, req.GivenDate, req.Amount, req.PaymentType, req.StudentId)
}
func (s *StudentService) ChangeUserBalanceHistoryByDebit(ctx context.Context, req *pb.ChangeUserBalanceHistoryByDebitRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Aborted, "error while getting company from context")
	}
	return s.repo.ChangeUserBalanceHistoryByDebit(companyId, req.StudentId, req.OldDebit, req.CurrentDebit, req.GivenDate, req.Comment, req.PaymentType, req.CreatedBy, req.CreatedByName, req.GroupId)
}
func (s *StudentService) CalculateDiscountSumma(ctx context.Context, req *pb.CalculateDiscountSummaRequest) (*pb.CalculateDiscountResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.Canceled, "error while getting company from context")
	}
	return s.repo.CalculateDiscountSumma(companyId, req.GroupId, req.StartDate, req.EndDate, req.DiscountPrice, req.StudentId, req.PaymentDate, req.StudentActivationDateInThisGroupWhilePayment)
}
