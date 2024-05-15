package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudentAPIHandler struct {
	studentRepo repository.StudentRepository
}

func NewStudentAPIHandler(studentRepo repository.StudentRepository) *StudentAPIHandler {
	return &StudentAPIHandler{studentRepo}
}

func (s *StudentAPIHandler) GetStudentByID(c *gin.Context) {
	// Implementation
}

func (s *StudentAPIHandler) AddStudent(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.studentRepo.Store(&student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "add student success"})
}
