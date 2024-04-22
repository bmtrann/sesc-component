package enrolment

import (
	"log"
	"net/http"

	"github.com/bmtrann/sesc-component/config"
	courseModel "github.com/bmtrann/sesc-component/internal/model/course"
	studentModel "github.com/bmtrann/sesc-component/internal/model/student"
	"github.com/bmtrann/sesc-component/internal/service"
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
	Student        studentModel.Student
	ServiceMessage string
}

type Payload struct {
	AccountId  string `json:"accountId"`
	StudentId  string `json:"studentId"`
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

	course, courseView, courseErr := handler.courseRepo.FindCourse(c.Request.Context(), payload.CourseName)

	if courseErr != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	studentId := payload.StudentId
	response := new(EnrolResponse)

	if studentId == "" {
		id := uuid.New()
		studentId = id.String()[:8]

		record := studentModel.MongoStudent{
			AccountId: payload.AccountId,
			StudentId: studentId,
			Courses:   []courseModel.Course{*course},
		}

		result, dbErr := handler.studentRepo.CreateStudent(c.Request.Context(), &record)
		response.Student = *result

		if dbErr != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := createFinanceAccount(studentId); err != nil {
			response.ServiceMessage = err.Error()
			c.JSON(http.StatusServiceUnavailable, response)
			return
		}

		if err := createLibraryAccount(studentId); err != nil {
			// Not throw since Create Invoice is more important
			log.Println(err)
			response.ServiceMessage = err.Error()
		}
	} else {
		err := handler.studentRepo.AddCourseToStudent(c.Request.Context(), studentId, course)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if err := createInvoice(studentId, courseView.Fees); err != nil {
		response.ServiceMessage = err.Error()
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusAccepted, response)
}

func createFinanceAccount(studentId string) error {
	return service.CreateFinanceAccount(studentId)
}

func createInvoice(studentId string, fees float32) error {
	return service.CreateInvoice(studentId, fees)
}

func createLibraryAccount(studentId string) error {
	return service.CreateLibraryAccount(studentId)
}
