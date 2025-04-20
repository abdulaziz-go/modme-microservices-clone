package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/etc"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func AddSmsCount(ctx *gin.Context) {

}

func DeleteSmsCount(ctx *gin.Context) {

}

func GetSmsTemplates(ctx *gin.Context) {

}

func SetSmsTemplates(ctx *gin.Context) {

}

func GetSmsTransactionDetail(ctx *gin.Context) {

}

func SendSmsDirectly(ctx *gin.Context) {

}
