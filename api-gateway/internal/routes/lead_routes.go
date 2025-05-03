package routes

import (
	client "api-gateway/internal/clients"
	"api-gateway/internal/etc"
	"api-gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

func LeadRoutes(api *gin.RouterGroup, userClient *client.UserClient) {
	lead := api.Group("/lead")
	{
		lead.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateLead)
		lead.POST("/get-lead-common", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetLeadCommon)
		lead.PUT("/update/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateLead)
		lead.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteLead)
		lead.GET("/get-all/:companyId", handlers.GetAllLead)
		lead.GET("/get-lead-reports", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetLeadReports)
	}
	expectation := api.Group("/expectation")
	{
		expectation.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateExpectation)
		expectation.PUT("/update/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateExpectation)
		expectation.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteExpectation)
	}
	set := api.Group("/set")
	{
		set.POST("/create", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.CreateSet)
		set.PUT("/update", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateSet)
		set.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteSet)
		set.PATCH("/change-to-group", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.ChangeToSet)
		set.GET("/get-by-id/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.GetByIdSet)
	}
	leadData := api.Group("/leadData")
	{
		leadData.POST("/create", handlers.CreateLeadData)
		leadData.PUT("/update", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.UpdateLeadData)
		leadData.DELETE("/delete/:id", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.DeleteLeadData)
		leadData.PATCH("/change-lead-data", etc.AuthMiddleware([]string{"ADMIN", "CEO"}, userClient), handlers.ChangeLeadData)
	}
}
