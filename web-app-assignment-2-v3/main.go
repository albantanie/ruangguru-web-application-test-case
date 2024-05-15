package main

import (
	"a21hc3NpZ25tZW50/api"
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/middleware"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"a21hc3NpZ25tZW50/service"
	"fmt"
	"time"

	_ "embed"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type APIHandler struct {
	UserAPIHandler     api.UserAPI
	CategoryAPIHandler api.CategoryAPI
	TaskAPIHandler     api.TaskAPI
}

func main() {
	gin.SetMode(gin.ReleaseMode) //release

	router := gin.New()
	db := db.NewDB()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s\"\n",
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	dbCredential := model.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "postgres",
		DatabaseName: "kampusmerdeka",
		Port:         5432,
		Schema:       "public",
	}

	conn, err := db.Connect(&dbCredential)
	if err != nil {
		panic(err)
	}

	conn.AutoMigrate(&model.User{}, &model.Category{}, &model.Task{})

	filebasedDb, err := filebased.InitDB()

	if err != nil {
		panic(err)
	}

	router = RunServer(router, filebasedDb)

	fmt.Println("Server is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func RunServer(gin *gin.Engine, filebasedDb *filebased.Data) *gin.Engine {
	userRepo := repo.NewUserRepo(filebasedDb)
	categoryRepo := repo.NewCategoryRepo(filebasedDb)
	taskRepo := repo.NewTaskRepo(filebasedDb)

	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	taskService := service.NewTaskService(taskRepo)

	userAPIHandler := api.NewUserAPI(userService)
	categoryAPIHandler := api.NewCategoryAPI(categoryService)
	taskAPIHandler := api.NewTaskAPI(taskService)

	apiHandler := APIHandler{
		UserAPIHandler:     userAPIHandler,
		CategoryAPIHandler: categoryAPIHandler,
		TaskAPIHandler:     taskAPIHandler,
	}

	user := gin.Group("/user")
	{
		// TODO: answer here
		user.Use(middleware.Auth())
		user.GET("/tasks", apiHandler.UserAPIHandler.GetUserTaskCategory)
		user.POST("/register", apiHandler.UserAPIHandler.Register)
		user.POST("/login", apiHandler.UserAPIHandler.Login)
	}

	task := gin.Group("/task")
	{
		task.Use(middleware.Auth())
		task.POST("/add", apiHandler.TaskAPIHandler.AddTask)
		task.GET("/get/:id", apiHandler.TaskAPIHandler.GetTaskByID)
		task.PUT("/update/:id", apiHandler.TaskAPIHandler.UpdateTask)
		task.DELETE("/delete/:id", apiHandler.TaskAPIHandler.DeleteTask)
		task.GET("/list", apiHandler.TaskAPIHandler.GetTaskList)
		task.GET("/category/:id", apiHandler.TaskAPIHandler.GetTaskListByCategory)
		task.POST("/category/add", apiHandler.CategoryAPIHandler.AddCategory)
		// TODO: answer here
	}

	category := gin.Group("/category")
	{
		category.Use(middleware.Auth())
		category.POST("/add", apiHandler.CategoryAPIHandler.AddCategory)
		category.GET("/get/:id", apiHandler.CategoryAPIHandler.GetCategoryByID)
		category.GET("/list", apiHandler.CategoryAPIHandler.GetCategoryList)
		category.PUT("/update/:id", apiHandler.CategoryAPIHandler.UpdateCategory)
		category.DELETE("/delete/:id", apiHandler.CategoryAPIHandler.DeleteCategory)
		// TODO: answer here
	}

	return gin
}
