package handlers

import (
	"api-gateway/grpc/proto/pb"
	"api-gateway/internal/etc"
	"api-gateway/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// CreateRoom godoc
// @Summary ADMIN , CEO
// @Description Create a new room based on the provided request data
// @Tags rooms
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body pb.CreateRoomRequest true "Request to create a room"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/create [post]
func CreateRoom(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	var req pb.CreateRoomRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusConflict, err.Error())
		return
	}
	resp, err := educationClient.CreateRoom(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateRoom godoc
// @Summary ADMIN , CEO
// @Description Update the details of an existing room based on the provided request data
// @Tags rooms
// @Accept json
// @Produce json
// @Security Bearer
// @Param request query pb.AbsRoom true "Request to update room"
// @Success 200 {object} utils.AbsResponse
// @Failure 422 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/update [put]
func UpdateRoom(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	var req pb.AbsRoom
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	resp, err := educationClient.UpdateRoom(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteRoom godoc
// @Summary ADMIN , CEO
// @Description Delete a room by its ID
// @Tags rooms
// @Produce json
// @Security Bearer
// @Param id path string true "Room ID"
// @Success 200 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/delete/{id} [delete]
func DeleteRoom(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.DeleteRoom(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllRoom godoc
// @Summary ADMIN , CEO
// @Description Retrieve all rooms
// @Tags rooms
// @Produce json
// @Security Bearer
// @Success 200 {object} pb.GetUpdateRoomAbs
// @Failure 500 {object} utils.AbsResponse
// @Router /api/room/get-all [get]
func GetAllRoom(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	rooms, err := educationClient.GetRoom(ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rooms)
	return
}

// CreateCourse godoc
// @Summary ADMIN , CEO
// @Description Create a new course based on the provided request data
// @Tags courses
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body pb.CreateCourseRequest true "Request to create a course"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/create [post]
func CreateCourse(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	var req pb.CreateCourseRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.CreateCourse(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateCourse godoc
// @Summary ADMIN , CEO
// @Description Update the details of an existing course based on the provided request data
// @Tags courses
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body pb.AbsCourse true "Request to update course"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/update [put]
func UpdateCourse(ctx *gin.Context) {
	var req pb.AbsCourse
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	resp, err := educationClient.UpdateCourse(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteCourse godoc
// @Summary ADMIN , CEO
// @Description Delete a course by its ID
// @Tags courses
// @Produce json
// @Security Bearer
// @Param id path string true "Course ID"
// @Success 200 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/delete/{id} [delete]
func DeleteCourse(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.DeleteCourse(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllCourse godoc
// @Summary ALL
// @Description Retrieve all courses
// @Tags courses
// @Produce json
// @Security Bearer
// @Success 200 {object} pb.GetUpdateCourseAbs
// @Failure 500 {object} utils.AbsResponse
// @Router /api/course/get-all [get]
func GetAllCourse(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	rooms, err := educationClient.GetCourse(ctxR)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, rooms)
	return
}

// GetCourseById godoc
// @Summary ADMIN , CEO
// @Description Retrieves a course by its ID for admin users.
// @Tags courses
// @Produce json
// @Security Bearer
// @Param id path string true "Course ID"
// @Success 200 {object} pb.AbsCourse "Successful response with course details"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/course/get-by-id/{id} [get]
func GetCourseById(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.GetCourseById(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// CreateGroup godoc
// @Summary ADMIN , CEO
// @Description Create a new group with provided details.
// @Tags groups
// @Accept json
// @Produce json
// @Security Bearer
// @Param group body pb.CreateGroupRequest true "Group Data"
// @Success 200 {object} utils.AbsResponse "Group successfully created"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/create [post]
func CreateGroup(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	var req pb.CreateGroupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.CreateGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// UpdateGroup godoc
// @Summary ADMIN , CEO
// @Description Update details of an existing group.
// @Tags groups
// @Accept json
// @Produce json
// @Security Bearer
// @Param group body pb.GetUpdateGroupAbs true "Group Data"
// @Success 200 {object} utils.AbsResponse "Group successfully updated"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/update [put]
func UpdateGroup(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	var req *pb.GetUpdateGroupAbs
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.UpdateGroup(ctxR, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteGroup godoc
// @Summary ADMIN , CEO
// @Description Delete a group by its ID.
// @Tags groups
// @Produce json
// @Security Bearer
// @Param id path string true "Group ID"
// @Success 200 {object} utils.AbsResponse "Group successfully deleted"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/delete/{id} [delete]
func DeleteGroup(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	resp, err := educationClient.DeleteGroup(ctxR, id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAllGroup godoc
// @Summary ADMIN , CEO
// @Description Retrieve a list of all groups.
// @Tags groups
// @Produce json
// @Security Bearer
// @Param group body pb.GetGroupsRequest true "Group Data"
// @Success 200 {array} pb.GetGroupsResponse "List of groups"
// @Failure 400 {object} utils.AbsResponse "Bad Request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/get-all [post]
func GetAllGroup(ctx *gin.Context) {
	var (
		req *pb.GetGroupsRequest
	)
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.GetAllGroup(ctxR, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
}

// GetGroupById godoc
// @Summary ADMIN
// @Description Retrieve details of a group by its ID.
// @Tags groups
// @Produce json
// @Security Bearer
// @Param id path string true "Group ID"
// @Success 200 {object} pb.GetGroupAbsResponse "Group details"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/get-by-id/{id} [get]
func GetGroupById(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	resp, err := educationClient.GetGroupById(ctxR, id, user.Id, user.Role)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
	return
}

// SetAttendance godoc
// @Summary TEACHER
// @Description Record attendance for a student in a group on a specific date.
// @Tags attendance
// @Produce json
// @Security Bearer
// @Param attendance body pb.SetAttendanceRequest true "Attendance details"
// @Success 200 {object} utils.AbsResponse "Attendance recorded successfully"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/attendance/set [post]
func SetAttendance(ctx *gin.Context) {
	var req pb.SetAttendanceRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	req.ActionById = user.Id
	req.ActionByRole = user.Role
	ctxR, cancelFunc := etc.NewTimoutContext(ctx)
	defer cancelFunc()
	resp, err := educationClient.SetAttendanceByGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetAttendance godoc
// @Summary ADMIN , TEACHER
// @Description Retrieve attendance records for students in a group over a specified date range.
// @Tags attendance
// @Produce json
// @Security Bearer
// @Param attendance body pb.GetAttendanceRequest true "Group ID and date range"
// @Success 200 {object} pb.GetAttendanceResponse "Attendance records"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/attendance/get-attendance [post]
func GetAttendance(ctx *gin.Context) {
	var req pb.GetAttendanceRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req.ActionRole = user.Role
	req.ActionId = user.Id
	resp, err := educationClient.GetAttendanceByGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
	return
}

// GetGroupByCourseId godoc
// @Summary ADMIN, TEACHER
// @Description Retrieve groups associated with a specific course ID.
// @Tags groups
// @Produce json
// @Security Bearer
// @Param courseId path string true "Course ID"
// @Success 200 {object} pb.GetGroupsByCourseResponse "Group details"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/get-by-course/{courseId} [get]
func GetGroupByCourseId(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	courseId := ctx.Param("courseId")
	resp, err := educationClient.GetGroupByCourseId(ctxR, courseId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
	return
}

// GetAllStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security Bearer
// @Param filter body pb.GetAllStudentRequest true "filter students"
// @Success 200 {object} pb.GetAllStudentResponse "List of students"
// @Failure 400 {object} utils.AbsResponse "Invalid condition"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/get-all [post]
func GetAllStudent(ctx *gin.Context) {
	var (
		req *pb.GetAllStudentRequest
	)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	response, err := educationClient.GetAllStudent(ctxR, req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// CreateStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security Bearer
// @Param request body pb.CreateStudentRequest true "Student details"
// @Success 200 {object} utils.AbsResponse "Created student details"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/create [post]
func CreateStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.CreateStudentRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.CreatedBy = "c1d6503f-31dc-4f99-b61f-2e4ebc7a7639"
	response, err := educationClient.CreateStudent(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// AddStudentToGroup godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security Bearer
// @Param request body pb.AddToGroupRequest true "Add student to group details"
// @Success 200 {object} utils.AbsResponse "Success message"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/add-to-group [post]
func AddStudentToGroup(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.AddToGroupRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	req.CreatedBy = user.Id
	response, err := educationClient.AddStudentToGroup(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, response.Status, response.Message)
	return
}

// UpdateStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security Bearer
// @Param request body pb.UpdateStudentRequest true "Updated student details"
// @Success 200 {object} utils.AbsResponse "Update success message"
// @Failure 400 {object} utils.AbsResponse "Invalid request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/update [put]
func UpdateStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.UpdateStudentRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response, err := educationClient.UpdateStudent(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, response.Status, response.Message)
	return
}

// DeleteStudent godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security Bearer
// @Param id path string true "Student ID"
// @Param returnMoney query bool true "Flag to determine if money should be returned"
// @Success 200 {object} utils.AbsResponse "Delete success message"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/delete/{id} [delete]
func DeleteStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	id := ctx.Param("id")
	returnMoney, err := strconv.ParseBool(ctx.Query("returnMoney"))
	if err != nil {
		return
	}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	response, err := educationClient.DeleteStudent(ctxR, id, returnMoney, user.Id, user.Name)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, response.Status, response.Message)
	return
}

// GetStudentById godoc
// @Summary ADMIN
// @Tags students
// @Produce json
// @Security Bearer
// @Param studentId path string true "Student ID"
// @Success 200 {object} pb.GetAllStudentResponse "List of students"
// @Failure 400 {object} utils.AbsResponse "Invalid condition"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/student/get-student-by-id/{studentId} [get]
func GetStudentById(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	studentId := ctx.Param("studentId")
	response, err := educationClient.GetStudentById(ctxR, studentId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// GetNotesByStudent godoc
// @Summary ADMIN
// @Description Get all notes associated with a specific student
// @Tags notes
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Success 200 {object} pb.GetNotesByStudent
// @Failure 500 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/student/note/get-notes/{studentId} [get]
func GetNotesByStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	studentId := ctx.Param("studentId")
	response, err := educationClient.GetNotesByStudentId(ctxR, studentId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// CreateNoteForStudent godoc
// @Summary ADMIN
// @Description Create a new note associated with a specific student
// @Tags notes
// @Accept json
// @Produce json
// @Param request body pb.CreateNoteRequest true "Note details"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/student/note/create [post]
func CreateNoteForStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.CreateNoteRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.CreateNoteForStudent(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// DeleteStudentNote godoc
// @Summary ADMIN
// @Description Delete a specific note associated with a student
// @Tags notes
// @Accept json
// @Produce json
// @Param noteId path string true "Note ID"
// @Success 200 {object} utils.AbsResponse
// @Failure 500 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/student/note/delete/{noteId} [delete]
func DeleteStudentNote(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	note := ctx.Param("noteId")
	resp, err := educationClient.DeleteNote(ctxR, note)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// SearchStudent godoc
// @Summary ADMIN
// @Description Search for students by phone number or name
// @Tags students
// @Accept json
// @Produce json
// @Param value path string true "Search value (phone number or name)"
// @Success 200 {object} pb.SearchStudentResponse
// @Failure 500 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/student/search-student/{value} [get]
func SearchStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	value := ctx.Param("value")
	resp, err := educationClient.SearchStudentByPhoneName(ctxR, value)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetHistoryGroup godoc
// @Summary ADMIN
// @Description Get the history of a specific group by its ID
// @Tags groups
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 200 {object} pb.GetHistoryGroupResponse
// @Failure 500 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/history/group/{groupId} [get]
func GetHistoryGroup(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	value := ctx.Param("groupId")
	resp, err := educationClient.GetHistoryGroupById(ctxR, value)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetHistoryStudent godoc
// @Summary ADMIN
// @Description Get the history of a specific student by their ID
// @Tags students
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Success 200 {object} pb.GetHistoryStudentResponse
// @Failure 500 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/history/student/{studentId} [get]
func GetHistoryStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	value := ctx.Param("studentId")
	resp, err := educationClient.GetHistoryStudentById(ctxR, value)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// TransferLessonDate transfers the lesson date for a specific course.
// @Summary ADMIN
// @Description Transfers the lesson date for a course
// @Tags lesson
// @Accept json
// @Produce json
// @Param request body pb.TransferLessonRequest true "Transfer Lesson Request"
// @Success 200 {object} utils.AbsResponse "Status and message"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/group/transfer-date [post]
func TransferLessonDate(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.TransferLessonRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.TransferLessonDate(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// ChangeConditionStudent changes the condition of a student.
// @Summary ADMIN, CEO
// @Description Changes the condition of a student based on provided details
// @Tags students
// @Accept json
// @Produce json
// @Param request body pb.ChangeConditionStudentRequest true "Change Condition Student Request"
// @Success 200 {object} utils.AbsResponse "Status and message"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Security Bearer
// @Router /api/student/change-condition [put]
func ChangeConditionStudent(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	req := pb.ChangeConditionStudentRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	req.ActionById = user.Id
	req.ActionByName = user.Name
	resp, err := educationClient.ChangeConditionStudent(ctxR, &req)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
	return
}

// GetInformationByTeacher godoc
// @Summary ADMIN , TEACHER , CEO
// @Description Get information about a specific teacher by their ID, with an option to filter archived data.
// @Tags groups
// @Param teacherId path string true "Teacher ID"
// @Param isArchived query bool true "Whether to include archived information"
// @Produce json
// @Success 200 {object} pb.GetGroupsByTeacherResponse
// @Failure 400 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/group/get-by-teacher/{teacherId} [get]
func GetInformationByTeacher(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	teacherId := ctx.Param("teacherId")
	isArchived, err := strconv.ParseBool(ctx.Query("isArchived"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.GetInformationByTeacher(ctxR, teacherId, isArchived)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetCommonInformationCompany godoc
// @Summary ADMIN , CEO
// @Description Get common information about company
// @Tags education
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 400 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/common-information-company [get]
func GetCommonInformationCompany(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	activeLeadCount := leadClient.GetActiveLeadCount(ctxR)
	activeStudentCount, activeGroupCount, leaveGroupCount, commonDebtorsCount, eleminatedInTrial := educationClient.GetCommonEducationInformation(ctxR)
	_, payInCurrentMonth := financeClient.GetCommonFinanceInformation(ctxR)
	response := make(map[string]int)
	response["activeLeadCount"] = activeLeadCount
	response["activeStudentsCount"] = activeStudentCount
	response["activeGroupCount"] = activeGroupCount
	response["debtorsCount"] = commonDebtorsCount
	response["payInCurrentMonth"] = payInCurrentMonth
	response["leaveGroupCount"] = leaveGroupCount
	response["eleminatedInTrial"] = eleminatedInTrial
	ctx.JSON(http.StatusOK, response)
	return
}

// GetChartIncome godoc
// @Summary CEO
// @Description Get information about a income
// @Tags education
// @Param from query string true "from"
// @Param to query string true "to"
// @Produce json
// @Success 200 {object} pb.GetCommonInformationResponse
// @Failure 400 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/get-chart-income [get]
func GetChartIncome(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()
	resp, err := financeClient.GetChartIncome(ctxR, ctx.Query("from"), ctx.Query("to"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
	return
}

// GetTableGroups godoc
// @Summary ADMIN , CEO, TEACHER
// @Description Get common information about company
// @Tags education
// @Produce json
// @Param dateType query string true "dateType"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.AbsResponse
// @Security Bearer
// @Router /api/get-table-groups [get]
func GetTableGroups(ctx *gin.Context) {
	ctxR, cancel := etc.NewTimoutContext(ctx)
	defer cancel()

	queryDateType := ctx.Query("dateType")

	response, err := educationClient.GetAllGroup(ctxR, &pb.GetGroupsRequest{
		IsArchived: false,
		Page: &pb.PageRequest{
			Page: 1,
			Size: 100000,
		},
	})
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var filteredGroups []*pb.GetGroupAbsResponse
	for _, group := range response.Groups {
		if group.DateType == queryDateType {
			filteredGroups = append(filteredGroups, group)
		}
	}

	finalResponse := map[string]interface{}{
		"groups": filteredGroups,
	}

	ctx.JSON(http.StatusOK, finalResponse)
}

// LeftAfterTrial
// @Summary ADMIN , CEO
// @Description Retrieve the data left after the trial period based on the provided from date, to date, page, and page size
// @Tags education
// @Accept json
// @Produce json
// @Param from path string true "Start date of the period"
// @Param to path string true "End date of the period"
// @Param page query string false "Page number" default(1)
// @Param page_size query string false "Page size" default(10)
// @Success 200 {object} pb.GetLeftAfterTrialPeriodResponse
// @Failure 400 {object} map[string]interface{}
// @Security Bearer
// @Router /api/group/left-after-trial/{from}/{to} [get]
func LeftAfterTrial(ctx *gin.Context) {
	from := ctx.Param("from")
	to := ctx.Param("to")
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	ctxR, cancelFunc := etc.NewTimoutContext(ctx)
	defer cancelFunc()
	resp, err := educationClient.GetLeftAfterTrialPeriod(ctxR, from, to, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// CompanyCreate
// @Summary SUPER_CEO
// @Description Create a new company
// @Tags company
// @Accept json
// @Produce json
// @Param request body pb.CreateCompanyRequest true "Company creation request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Router /api/company/create [post]
// @Security Bearer
func CompanyCreate(ctx *gin.Context) {
	var req pb.CreateCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.CreateCompanyRequest(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
}

// GetAllCompanies
// @Summary SUPER_CEO
// @Description Retrieve a paginated list of all companies
// @Tags company
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param filter query string false "filter only accept (active , no_active , demo)"
// @Success 200 {object} pb.GetAllResponse
// @Failure 400 {object} utils.AbsResponse
// @Router /api/company/get-all [get]
// @Security Bearer
func GetAllCompanies(ctx *gin.Context) {
	resp, err := educationClient.GetAllCompanies(ctx.Query("page"), ctx.Query("size"), ctx.Query("filter"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetOneCompany
func GetOneCompany(ctx *gin.Context) {

}

// CompanyUpdate
// @Summary SUPER_CEO
// @Description Update an existing company by ID
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param request body pb.UpdateCompanyRequest true "Company update request"
// @Success 200 {object} utils.AbsResponse
// @Failure 400 {object} utils.AbsResponse
// @Router /api/company/update [put]
// @Security Bearer
func CompanyUpdate(ctx *gin.Context) {
	var req pb.UpdateCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.UpdateCompany(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, resp.Status, resp.Message)
}

// GetCompanyBySubdomain
// @Summary ALL
// @Description Retrieve the data company by domain
// @Tags company
// @Accept json
// @Produce json
// @Param domain path string true "Start date of the period"
// @Success 200 {object} pb.GetCompanyResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/company/subdomain/{domain} [get]
func GetCompanyBySubdomain(ctx *gin.Context) {
	domain := ctx.Param("domain")
	resp, err := educationClient.GetCompanyBySubdomain(domain)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// UploadImage
// @Summary Upload an image
// @Description Upload an image file and save it to the server
// @Tags image
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file to upload"
// @Success 200 {object} utils.AbsResponse "Success response with file URL"
// @Failure 400 {object} utils.AbsResponse "Bad request or file error"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/image/upload [post]
func UploadImage(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image file: " + err.Error()})
		return
	}
	defer file.Close()
	const maxFileSize = 5 << 20
	if header.Size > maxFileSize {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is too large. Maximum size is 5MB"})
		return
	}
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("image_%d%s", time.Now().UnixNano(), ext)

	uploadDir := "/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory: " + err.Error()})
		return
	}

	filePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
		return
	}
	defer out.Close()

	if _, err := file.Seek(0, 0); err == nil {
		if _, err := out.ReadFrom(file); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write image: " + err.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset file pointer: " + err.Error()})
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, fmt.Sprintf("https://backend.livesphere.uz/api/image/get-image?filename=%v", filename))
}

// GetImage
// @Summary Get an uploaded image
// @Description Retrieve an uploaded image by filename
// @Tags image
// @Accept json
// @Produce image/jpeg
// @Produce image/png
// @Produce image/gif
// @Param filename query string true "Filename of the image to retrieve"
// @Success 200 "Image file"
// @Failure 400 {object} utils.AbsResponse "Bad request, filename missing"
// @Failure 404 {object} utils.AbsResponse "Image not found"
// @Router /api/image/get-image [get]
func GetImage(ctx *gin.Context) {
	filename := ctx.Query("filename")
	if filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}
	uploadDir := "/uploads"
	filePath := filepath.Join(uploadDir, filename)
	if !fileExists(filePath) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	ctx.File(filePath)
}

func fileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

// TariffCreate
// @Summary SUPER_CEO
// @Description Create a new tariff with the provided details
// @Tags tariff
// @Accept json
// @Produce json
// @Param tariff body pb.Tariff true "Tariff data"
// @Success 200 {object} utils.AbsResponse "created"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Router /api/company/tariff/create [post]
// @Security Bearer
func TariffCreate(ctx *gin.Context) {
	var tariff pb.Tariff
	if err := ctx.ShouldBindJSON(&tariff); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	_, err := educationClient.CreateTariff(&tariff)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, 200, "created")
}

// TariffGetAll
// @Summary SUPER_CEO
// @Description Retrieve a list of all tariffs
// @Tags tariff
// @Accept json
// @Produce json
// @Success 200 {object} pb.TariffList "List of tariffs"
// @Failure 500 {object} utils.AbsResponse "Internal server error"
// @Router /api/company/tariff/get-all [get]
// @Security Bearer
func TariffGetAll(ctx *gin.Context) {
	resp := educationClient.GetAllTariff()
	ctx.JSON(http.StatusOK, resp)
}

// TariffUpdate
// @Summary SUPER_CEO
// @Description Update the details of an existing tariff
// @Tags tariff
// @Accept json
// @Produce json
// @Param tariff body pb.Tariff true "Updated tariff data"
// @Success 200 {object} utils.AbsResponse "updated"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Router /api/company/tariff/update [put]
// @Security Bearer
func TariffUpdate(ctx *gin.Context) {
	var tariff pb.Tariff
	if err := ctx.ShouldBindJSON(&tariff); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	_, err := educationClient.UpdateTariff(&tariff)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, 200, "updated")
}

// TariffDelete
// @Summary Delete a tariff
// @Description Delete a tariff by its ID
// @Tags tariff
// @Accept json
// @Produce json
// @Param id path int true "ID of the tariff to delete"
// @Success 200 {object} utils.AbsResponse "deleted"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Failure 404 {object} utils.AbsResponse "Tariff not found"
// @Router /api/company/tariff/delete/{id} [delete]
// @Security Bearer
func TariffDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	idn, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err = educationClient.DeleteTariff(int32(idn))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, 200, "deleted")
}

// FinanceCreate godoc
// @Summary Create a finance record
// @Description Create a new company finance record
// @Tags companyFinance
// @Accept json
// @Produce json
// @Param body body pb.CompanyFinance true "Company finance details"
// @Success 200 {object} utils.AbsResponse "created"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Router /api/company/finance/create [post]
// @Security Bearer
func FinanceCreate(ctx *gin.Context) {
	req := pb.CompanyFinance{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err := educationClient.FinanceCreate(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, 200, "created")
}

// FinanceDelete godoc
// @Summary Delete a finance record
// @Description Delete a company finance record by its ID
// @Tags companyFinance
// @Accept json
// @Produce json
// @Param id query string true "ID of the company finance record to delete"
// @Success 200 {object} utils.AbsResponse "deleted"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Router /api/company/finance/delete [delete]
// @Security Bearer
func FinanceDelete(ctx *gin.Context) {
	financeCompanyId := ctx.Query("id")
	_, err := educationClient.FinanceDelete(financeCompanyId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, "deleted")
}

// FinanceGetAll godoc
// @Summary Get all finance records
// @Description Retrieve a list of all company finance records
// @Tags companyFinance
// @Accept json
// @Produce json
// @Param body body pb.PageRequest true "Pagination details"
// @Success 200 {object} pb.CompanyFinanceList "list of finance records"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Router /api/company/finance/get-all [post]
// @Security Bearer
func FinanceGetAll(ctx *gin.Context) {
	req := pb.PageRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.FinanceGetAll(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
}

// FinanceGetByCompany godoc
// @Summary Get finance records by company
// @Description Retrieve finance records filtered by company
// @Tags companyFinance
// @Accept json
// @Produce json
// @Param body body pb.PageRequest true "Pagination and filtering details"
// @Success 200 {object} pb.CompanyFinanceSelfList "list of finance records for a specific company"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Router /api/company/finance/get-by-company [post]
// @Security Bearer
func FinanceGetByCompany(ctx *gin.Context) {
	req := pb.PageRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := educationClient.FinanceGetByCompany(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, &resp)
}

// FinanceUpdateByCompany godoc
// @Summary Update finance records by company
// @Description Update finance details for a specific company
// @Tags companyFinance
// @Accept json
// @Produce json
// @Param body body pb.CompanyFinance true "Finance details to be updated"
// @Success 200 {object} pb.CompanyFinance "Updated finance details"
// @Failure 400 {object} utils.AbsResponse "Bad request"
// @Router /api/company/finance/update [put]
// @Security Bearer
func FinanceUpdateByCompany(ctx *gin.Context) {
	req := pb.CompanyFinance{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	res, err := educationClient.FinanceUpdate(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// GetStatisticCompany godoc
// @Summary SUPER_CEO
// @Description
// @Tags company
// @Accept json
// @Produce json
// @Param body body pb.GetStatisticRequest true "Параметры запроса статистики компании"
// @Success 200 {object} pb.GetStatisticResponse "Успешный ответ со статистикой"
// @Failure 400 {object} utils.AbsResponse "Неверный запрос"
// @Router /api/company/get-statistic [post]
// @Security Bearer
func GetStatisticCompany(ctx *gin.Context) {
	req := pb.GetStatisticRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	response, err := educationClient.GetStatisticCompany(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
}
