package api

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"

	"net/http"

	"github.com/gin-gonic/gin"
)

type CourseAPI interface {
	AddCourse(c *gin.Context)
}

type courseAPI struct {
	courseRepo repo.CourseRepository
}

func NewCourseAPI(courseRepo repo.CourseRepository) *courseAPI {
	return &courseAPI{courseRepo}
}

func (cr *courseAPI) AddCourse(c *gin.Context) {
	course := &model.Course{} // Create a pointer to a model.Course

	// Attempt to store the course in the repository
	if err := cr.courseRepo.Store(course); err != nil {
		// If an error occurs, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add course"})
		return
	}

	// If the course is successfully stored, return a success response
	c.JSON(http.StatusOK, "add")
}
