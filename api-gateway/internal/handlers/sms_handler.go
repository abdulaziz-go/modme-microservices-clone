package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/etc"
	"api-gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSmsLogs godoc
// @Summary Get SMS logs
// @Description ADMIN, CEO, FINANCIST, SUPER_CEO roles can retrieve SMS logs
// @Tags company-sms
// @Accept json
// @Produce json
// @Param body body pb.GetSmsLogRequest true "Request for fetching SMS logs"
// @Success 200 {object} pb.GetSmsLogResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/company/sms/get-logs [post]
// @Security Bearer
func GetSmsLogs(ctx *gin.Context) {
	var (
		req pb.GetSmsLogRequest
	)
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := educationClient.GetSmsLogs(ctxR, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// AddSmsCount godoc
// @Summary Add SMS count to a company
// @Description SUPER_CEO role can add SMS count
// @Tags company-sms
// @Accept json
// @Produce json
// @Param body body pb.AddSmsRequest true "Add SMS count request"
// @Success 200 {object} pb.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/company/sms/add-sms-count [post]
// @Security Bearer
func AddSmsCount(ctx *gin.Context) {
	var (
		req pb.AddSmsRequest
	)
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := educationClient.AddSms(ctxR, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// DeleteSmsCount godoc
// @Summary Delete SMS count by transaction ID
// @Description SUPER_CEO role can delete SMS transaction
// @Tags company-sms
// @Accept json
// @Produce json
// @Param transactionId query string true "Transaction ID"
// @Success 200 {object} pb.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/company/sms/delete-sms-count [delete]
// @Security Bearer
func DeleteSmsCount(ctx *gin.Context) {
	id := ctx.Query("transactionId")
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	response, err := educationClient.DeleteSms(ctxR, &pb.DeleteAbsRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetSmsTemplates godoc
// @Summary Get SMS templates
// @Description ADMIN, CEO roles can fetch templates
// @Tags company-sms-template
// @Accept json
// @Produce json
// @Param body body pb.GetSmsTemplateRequest true "Get templates request"
// @Success 200 {object} pb.GetSmsTemplateResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/company/sms/template/get-templates [post]
// @Security Bearer
func GetSmsTemplates(ctx *gin.Context) {
	var (
		req pb.GetSmsTemplateRequest
	)
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := educationClient.GetSmsTemplate(ctxR, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// SetSmsTemplates godoc
// @Summary Set or update SMS templates
// @Description CEO, ADMIN roles can update templates
// @Tags company-sms-template
// @Accept json
// @Produce json
// @Param request body pb.SetSmsTemplateRequest true "Set template request"
// @Success 200 {object} pb.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/company/sms/template/set-template [put]
// @Security Bearer
func SetSmsTemplates(ctx *gin.Context) {
	var (
		req pb.SetSmsTemplateRequest
	)
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := educationClient.SetSmsTemplate(ctxR, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetSmsTransactionDetail godoc
// @Summary Get all SMS transactions (paginated)
// @Description ADMIN, CEO, SUPER_CEO roles can access this
// @Tags company-sms
// @Accept json
// @Produce json
// @Param body body pb.PageRequest true "Pagination request"
// @Success 200 {object} pb.GetSmsTransactionDetailResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/company/sms/get-sms-count-all [post]
// @Security Bearer
func GetSmsTransactionDetail(ctx *gin.Context) {
	var (
		req pb.PageRequest
	)
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := educationClient.GetSmsTransactionDetail(ctxR, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// SendSmsDirectly godoc
// @Summary Send an SMS message directly
// @Description ADMIN, CEO roles can send SMS
// @Tags company-sms
// @Accept json
// @Produce json
// @Param request body pb.SendSmsDirectlyRequest true "Send SMS directly request"
// @Success 200 {object} pb.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/company/sms/send-directly [post]
// @Security Bearer
func SendSmsDirectly(ctx *gin.Context) {
	var (
		req pb.SendSmsDirectlyRequest
	)
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userDetail, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.CreatorName = userDetail.Name
	req.CreatorId = userDetail.Id
	response, err := educationClient.SendSmsDirectly(ctxR, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}
