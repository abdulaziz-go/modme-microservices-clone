package service

import (
	"context"
	"education-service/internal/repository"
	"education-service/internal/utils"
	"education-service/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SmsService struct {
	pb.UnimplementedSmsServiceServer
	smsRepo *repository.SmsRepository
}

func NewSmsService(repo *repository.SmsRepository) *SmsService {
	return &SmsService{
		smsRepo: repo,
	}
}

func (s *SmsService) GetSmsLogs(ctx context.Context, req *pb.GetSmsLogRequest) (*pb.GetSmsLogResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.FailedPrecondition, "Company ID is required")
	}

	return s.smsRepo.GetSmsLog(companyId, req.StudentId, req.PageRequest)
}

func (s *SmsService) AddSms(ctx context.Context, req *pb.AddSmsRequest) (*pb.AbsResponse, error) {
	return s.smsRepo.AddSms(req)
}
func (s *SmsService) DeleteSms(ctx context.Context, req *pb.DeleteAbsRequest) (*pb.AbsResponse, error) {
	return s.smsRepo.DeleteSms(req.Id)
}
func (s *SmsService) GetSmsTransactionDetail(ctx context.Context, req *pb.PageRequest) (*pb.GetSmsTransactionDetailResponse, error) {
	return s.smsRepo.GetSmsTransactionDetail(req.Page, req.Size, req.CompanyId)
}
func (s *SmsService) GetSmsTemplate(ctx context.Context, req *pb.GetSmsTemplateRequest) (*pb.GetSmsTemplateResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.FailedPrecondition, "Company ID is required")
	}

	return s.smsRepo.GetSmsTemplate(req.SmsType, companyId)
}
func (s *SmsService) SetSmsTemplate(ctx context.Context, req *pb.SetSmsTemplateRequest) (*pb.AbsResponse, error) {
	return nil, nil
}
func (s *SmsService) SendSmsDirectly(ctx context.Context, req *pb.SendSmsDirectlyRequest) (*pb.AbsResponse, error) {
	companyId := utils.GetCompanyId(ctx)
	if companyId == "" {
		return nil, status.Error(codes.FailedPrecondition, "Company ID is required")
	}

	return s.smsRepo.SendSmsDirectly(req, companyId)
}
