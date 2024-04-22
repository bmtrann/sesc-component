package profile

import (
	"net/http"

	"github.com/bmtrann/sesc-component/internal/exception"
	studentModel "github.com/bmtrann/sesc-component/internal/model/student"
	"github.com/bmtrann/sesc-component/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileService interface {
	GetProfile()
	UpdateProfile()
	GetGraduationStatus()
}

type GetProfileResponse struct {
	Student studentModel.Student
}

type UpdateProfilePayload struct {
	StudentId string `json:"studentId"`
	FirstName string `json:"firstName"`
	Surname   string `json:"surname"`
}

type ProfileHandler struct {
	studentRepo *studentModel.StudentRepository
}

func InitProfileHandler(db *mongo.Database, collecion string) *ProfileHandler {
	return &ProfileHandler{
		studentRepo: studentModel.NewStudentRepository(db, collecion),
	}
}

func (handler *ProfileHandler) GetProfile(c *gin.Context) {
	student_id := c.Param("id")
	student, err := handler.studentRepo.GetStudent(c.Request.Context(), student_id)

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, GetProfileResponse{Student: *student})
}

func (handler *ProfileHandler) UpdateProfile(c *gin.Context) {
	payload := new(UpdateProfilePayload)

	if err := c.BindJSON(payload); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	data := make(map[string]string)

	if payload.FirstName != "" {
		data["firstName"] = payload.FirstName
	}

	if payload.Surname != "" {
		data["surname"] = payload.Surname
	}

	err := handler.studentRepo.UpdateStudentProfile(c.Request.Context(), payload.StudentId, data)
	if err != nil {
		if err == exception.ErrUserNotFound {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (handler *ProfileHandler) GetGraduationStatus(c *gin.Context) {
	studentId := c.Param("id")

	status, err := service.GetGraduationStatus(studentId)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]bool{
		"graduationStatus": !status,
	})
}
