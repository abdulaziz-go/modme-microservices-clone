package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/etc"
	"api-gateway/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateLead godoc
// @Summary ADMIN
// @Description Create a new lead with the given title.
// @Tags leads
// @Accept json
// @Produce json
// @Param title query string true "Title of the lead"
// @Success 200 {object} utils.AbsResponse "Lead successfully created"
// @Failure 409 {object} utils.AbsResponse "Lead creation conflict"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/lead/create [post]
func CreateLead(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	value := ctx.Query("title")
	resp, err := leadClient.CreateLead(ctxR, value)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetLeadCommon godoc
// @Summary ADMIN
// @Description Get a common lead by ID and type.
// @Tags leads
// @Accept json
// @Produce json
// @Param req body pb.GetLeadCommonRequest true "Lead ID"
// @Success 200 {object} pb.GetLeadCommonResponse "Lead details retrieved"
// @Failure 400 {object} utils.AbsResponse "Bad Request"
// @Failure 500 {object} utils.AbsResponse "Internal Server Error"
// @Security Bearer
// @Router /api/lead/get-lead-common [post]
func GetLeadCommon(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()

	var req pb.GetLeadCommonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := leadClient.GetLeadCommon(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// UpdateLead godoc
// @Summary ADMIN
// @Description Update an existing lead by ID.
// @Tags leads
// @Accept json
// @Produce json
// @Param id path string true "Lead ID"
// @Param title query string true "Title of the lead"
// @Success 200 {object} utils.AbsResponse "Lead updated successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/lead/update/{id} [put]
func UpdateLead(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	title := ctx.Query("title")
	resp, err := leadClient.UpdateLead(ctxR, id, title)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return

}

// DeleteLead godoc
// @Summary ADMIN
// @Description Delete a lead by ID.
// @Tags leads
// @Accept json
// @Produce json
// @Param id path string true "Lead ID"
// @Success 200 {object} utils.AbsResponse "Lead deleted successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/lead/delete/{id} [delete]
func DeleteLead(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := leadClient.DeleteLead(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// CreateExpectation godoc
// @Summary ADMIN
// @Description Create a new expectation.
// @Tags expectations
// @Accept json
// @Produce json
// @Param title query string true "Title of the expectation"
// @Success 200 {object} utils.AbsResponse "Expectation created successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/expectation/create [post]
func CreateExpectation(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	title := ctx.Query("title")
	resp, err := leadClient.CreateExpect(ctxR, title)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateExpectation godoc
// @Summary ADMIN
// @Description Update an existing expectation by ID.
// @Tags expectations
// @Accept json
// @Produce json
// @Param id path string true "Expectation ID"
// @Param title query string true "Title of the expectation"
// @Success 200 {object} utils.AbsResponse "Expectation updated successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/expectation/update/{id} [put]
func UpdateExpectation(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	title := ctx.Query("title")
	resp, err := leadClient.UpdateExpect(ctxR, id, title)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteExpectation godoc
// @Summary ADMIN
// @Description Delete an expectation by ID.
// @Tags expectations
// @Accept json
// @Produce json
// @Param id path string true "Expectation ID"
// @Success 200 {object} utils.AbsResponse "Expectation deleted successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/expectation/delete/{id} [delete]
func DeleteExpectation(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := leadClient.DeleteExpect(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// CreateSet godoc
// @Summary ADMIN
// @Description Create a new set.
// @Tags sets
// @Accept json
// @Produce json
// @Param data body pb.CreateSetRequest true "Set data"
// @Success 200 {object} utils.AbsResponse "Set created successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/set/create [post]
func CreateSet(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.CreateSetRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := leadClient.CreateSet(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateSet godoc
// @Summary ADMIN
// @Description Update an existing set by ID.
// @Tags sets
// @Accept json
// @Produce json
// @Param data body pb.UpdateSetRequest true "Set data"
// @Success 200 {object} utils.AbsResponse "Set updated successfully"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/set/update [put]
func UpdateSet(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.UpdateSetRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := leadClient.UpdateSet(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteSet godoc
// @Summary ADMIN
// @Description Delete a set by ID.
// @Tags sets
// @Accept json
// @Produce json
// @Param id path string true "Set ID"
// @Success 200 {object} utils.AbsResponse "Set deleted successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/set/delete/{id} [delete]
func DeleteSet(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := leadClient.DeleteSet(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// CreateLeadData godoc
// @Summary ADMIN
// @Description Create lead data.
// @Tags leadData
// @Accept json
// @Produce json
// @Param data body pb.CreateLeadDataRequest true "Lead data"
// @Success 200 {object} utils.AbsResponse "Lead data created successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/leadData/create [post]
func CreateLeadData(ctx *gin.Context) {
	req := pb.CreateLeadDataRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctxR, cancel := etc.NewTimeoutContextManual(ctx, req.CompanyId)
	defer cancel()
	resp, err := leadClient.CreateLeadData(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteLeadData godoc
// @Summary ADMIN
// @Description Delete lead data by ID.
// @Tags leadData
// @Accept json
// @Produce json
// @Param id path string true "Lead data ID"
// @Success 200 {object} utils.AbsResponse "Lead data deleted successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/leadData/delete/{id} [delete]
func DeleteLeadData(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := leadClient.DeleteLeadData(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateLeadData godoc
// @Summary ADMIN
// @Description Update lead data by ID.
// @Tags leadData
// @Accept json
// @Produce json
// @Param data body pb.UpdateLeadDataRequest true "Lead data"
// @Success 200 {object} utils.AbsResponse "Lead data updated successfully"
// @Failure 409 {object} utils.AbsResponse "Conflict occurred"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/leadData/update [put]
func UpdateLeadData(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.UpdateLeadDataRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := leadClient.UpdateLeadData(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// ChangeLeadData godoc
// @Summary Change lead data
// @Description Update the data associated with a lead
// @Tags leadData
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body pb.ChangeLeadPlaceRequest true "Lead change request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 401 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Router /api/leadData/change-lead-data [patch]
func ChangeLeadData(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.ChangeLeadPlaceRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := leadClient.ChangeLeadPlace(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllLead godoc
// @Summary ALL
// @Description Update the data associated with a lead
// @Tags leads
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} pb.GetLeadListResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 401 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Router /api/lead/get-all/{companyId} [get]
func GetAllLead(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimeoutContextManual(ctx, ctx.Param("companyId"))
	defer cancel()

	resp, err := leadClient.GetAllLead(ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// ChangeToSet godoc
// @Summary ADMIN
// @Description Change the lead set to a group based on the provided request data
// @Tags sets
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body pb.ChangeToSetRequest true "Request to change set to group"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 409 {object} utils.AbsResponse
// @Router /api/set/change-to-group [patch]
func ChangeToSet(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	var req pb.ChangeToSetRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := leadClient.ChangeSetToGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetLeadReports handles the retrieval of lead reports based on specified date range.
// @Summary CEO , ADMIN
// @Description Retrieves lead reports for a given date range.
// @Tags leads
// @Accept json
// @Produce json
// @Param from query string true "Start date in YYYY-MM-DD format"
// @Param till query string true "End date in YYYY-MM-DD format"
// @Success 200 {object} utils.AbsResponse "Lead reports response"
// @Failure 409 {object} utils.AbsResponse "Conflict error with details"
// @Security Bearer
// @Router /api/lead/get-lead-reports [get]
func GetLeadReports(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	from := ctx.Query("from")
	till := ctx.Query("till")
	resp, err := leadClient.GetLeadReports(from, till, ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetByIdSet retrieves lead details by ID.
// @Summary ADMIN , CEO
// @Description Retrieves lead details based on the provided ID.
// @Tags sets
// @Accept json
// @Produce json
// @Param id path string true "Lead ID"
// @Success 200 {object} pb.SetDataResponse "Successful response with lead details"
// @Failure 404 {object} utils.AbsResponse "Lead not found"
// @Failure 409 {object} utils.AbsResponse "Conflict error with details"
// @Security Bearer
// @Router /api/set/get-by-id/{id} [get]
func GetByIdSet(ctx *gin.Context) {
	ctxR, cancelFunc := etc.NewTimoutContext(ctx)
	defer cancelFunc()
	id := ctx.Param("id")
	resp, err := leadClient.GetByIdSet(id, ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
