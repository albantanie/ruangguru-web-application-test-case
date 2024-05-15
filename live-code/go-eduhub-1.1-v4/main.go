package main

import (
	"a21hc3NpZ25tZW50/api"
	repo "a21hc3NpZ25tZW50/repository"
	"fmt"
	"time"

	_ "embed"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	StudentAPIHandler api.StudentAPI
	CourseAPIHandler  api.CourseAPI
}

func main() {
	gin.SetMode(gin.ReleaseMode) //release

	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s\"\n",
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router = RunServer(router)

	fmt.Println("Server is running on port 8080")
	var err error
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func RunServer(gin *gin.Engine) *gin.Engine {
	studentRepo := repo.NewStudentRepo()
	courseRepo := repo.NewCourseRepo()

	studentAPIHandler := api.NewStudentAPI(studentRepo)
	courseAPIHandler := api.NewCourseAPI(courseRepo)

	apiHandler := APIHandler{
		StudentAPIHandler: studentAPIHandler,
		CourseAPIHandler:  courseAPIHandler,
	}

	student := gin.Group("/student")
	{
		student.POST("/add", apiHandler.StudentAPIHandler.AddStudent)
		student.GET("/gets", apiHandler.StudentAPIHandler.GetStudents)
		student.GET("/get/:id", apiHandler.StudentAPIHandler.GetStudentByID)
	}

	course := gin.Group("/course")
	{
		course.POST("/add", apiHandler.CourseAPIHandler.AddCourse)
	}

	return gin
}
