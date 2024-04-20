package enrolment

import (
	"net/http"

	"github.com/bmtrann/sesc-component/config"
	"github.com/bmtrann/sesc-component/internal/exception"
	courseModel "github.com/bmtrann/sesc-component/internal/model/course"
	studentModel "github.com/bmtrann/sesc-component/internal/model/student"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type EnrolmentService interface {
	List()
	View()
	Enrol()
}

type EnrolmentHandler struct {
	courseRepo  *courseModel.CourseRepository
	studentRepo *studentModel.StudentRepository
}

type ListResponse struct {
	Courses []courseModel.CourseView
}

type EnrolResponse struct {
	Student studentModel.Student
}

type Payload struct {
	AccountId  string `json:"accountId"`
	CourseName string `json:"courseName"`
}

func InitEnrolmentHandler(db *mongo.Database, dbConfig *config.DBConfig) *EnrolmentHandler {
	return &EnrolmentHandler{
		courseModel.NewCourseRepository(db, dbConfig.CourseCollection),
		studentModel.NewStudentRepository(db, dbConfig.StudentCollection),
	}
}

func (handler *EnrolmentHandler) List(c *gin.Context) {
	courses, err := handler.courseRepo.GetCourses(c.Request.Context(), nil)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ListResponse{Courses: courses})
}

func (handler *EnrolmentHandler) View(c *gin.Context) {
	studentId := c.Param("id")

	student, err := handler.studentRepo.GetStudent(c.Request.Context(), studentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	studentCourses := student.Courses
	courses, err := handler.courseRepo.GetCourses(c.Request.Context(), studentCourses)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ListResponse{Courses: courses})
}

func (handler *EnrolmentHandler) Enrol(c *gin.Context) {
	payload := new(Payload)

	if err := c.BindJSON(payload); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	course, _, courseErr := handler.courseRepo.FindCourse(c.Request.Context(), payload.CourseName)

	if courseErr != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err := handler.studentRepo.AddCourseToStudent(c.Request.Context(), payload.AccountId, course)

	if err != nil {
		if err == exception.ErrUserNotFound {
			id := uuid.New()

			record := studentModel.MongoStudent{
				AccountId: payload.AccountId,
				StudentId: id.String()[:8],
				Courses:   []courseModel.Course{*course},
			}

			student, _ := handler.studentRepo.CreateStudent(c.Request.Context(), &record)
			c.JSON(http.StatusOK, EnrolResponse{Student: *student})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
