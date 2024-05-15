package main_test

import (
	main "a21hc3NpZ25tZW50"
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/middleware"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"a21hc3NpZ25tZW50/service"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func SetCookie(mux *gin.Engine) *http.Cookie {
	login := model.UserLogin{
		Email:    "test@mail.com",
		Password: "testing123",
	}

	body, _ := json.Marshal(login)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/user/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(w, r)

	var cookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "session_token" {
			cookie = c
		}
	}

	return cookie
}

var _ = Describe("Task Tracker Plus", Ordered, func() {
	var apiServer *gin.Engine

	var userRepo repo.UserRepository
	var categoryRepo repo.CategoryRepository
	var taskRepo repo.TaskRepository

	var userService service.UserService
	var categoryService service.CategoryService
	var taskService service.TaskService

	var insertCategories []model.Category
	var insertTasks []model.Task

	var expectedUserTask []model.UserTaskCategory

	var err error

	var filebasedDb *filebased.Data

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode) //release

		os.Remove("file.db")
		filebasedDb = &filebased.Data{}

		filebasedDb, err = filebased.InitDB()

		Expect(err).ShouldNot(HaveOccurred())

		// _, err = db.Connect(&dbCredential)
		Expect(err).ShouldNot(HaveOccurred())

		userRepo = repo.NewUserRepo(filebasedDb)
		categoryRepo = repo.NewCategoryRepo(filebasedDb)
		taskRepo = repo.NewTaskRepo(filebasedDb)

		userService = service.NewUserService(userRepo)
		categoryService = service.NewCategoryService(categoryRepo)
		taskService = service.NewTaskService(taskRepo)

		Expect(err).ShouldNot(HaveOccurred())

		apiServer = gin.New()
		apiServer = main.RunServer(apiServer, filebasedDb)

		expectedUserTask = []model.UserTaskCategory{
			{
				ID:       1,
				Fullname: "test",
				Email:    "test@mail.com",
				Task:     "Task 2",
				Deadline: "2023-06-01",
				Priority: 1,
				Status:   "Completed",
				Category: "Category 2",
			},
			{
				ID:       1,
				Fullname: "test",
				Email:    "test@mail.com",
				Task:     "Task 5",
				Deadline: "2023-06-07",
				Priority: 5,
				Status:   "In Progress",
				Category: "Category 3",
			},
		}

		// Init test data:
		insertCategories = []model.Category{
			{ID: 1, Name: "Category 1"},
			{ID: 2, Name: "Category 2"},
			{ID: 3, Name: "Category 3"},
			{ID: 4, Name: "Category 4"},
			{ID: 5, Name: "Category 5"},
		}

		for _, v := range insertCategories {
			err := categoryRepo.Store(&v)
			Expect(err).ShouldNot(HaveOccurred())
		}

		insertTasks = []model.Task{
			{
				ID:         1,
				Title:      "Task 1",
				Deadline:   "2023-05-30",
				Priority:   2,
				Status:     "In Progress",
				CategoryID: 1,
				UserID:     2,
			},
			{
				ID:         2,
				Title:      "Task 2",
				Deadline:   "2023-06-01",
				Priority:   1,
				Status:     "Completed",
				CategoryID: 2,
				UserID:     1,
			},
			{
				ID:         3,
				Title:      "Task 3",
				Deadline:   "2023-06-02",
				Priority:   4,
				Status:     "Completed",
				CategoryID: 1,
				UserID:     3,
			},
			{
				ID:         4,
				Title:      "Task 4",
				Deadline:   "2023-06-02",
				Priority:   3,
				Status:     "Completed",
				CategoryID: 1,
				UserID:     4,
			},
			{
				ID:         5,
				Title:      "Task 5",
				Deadline:   "2023-06-07",
				Priority:   5,
				Status:     "In Progress",
				CategoryID: 3,
				UserID:     1,
			},
		}

		for _, v := range insertTasks {
			err := taskRepo.Store(&v)
			Expect(err).ShouldNot(HaveOccurred())
		}

		reqRegister := model.UserRegister{
			Fullname: "test",
			Email:    "test@mail.com",
			Password: "testing123",
		}

		reqBody, _ := json.Marshal(reqRegister)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")
		apiServer.ServeHTTP(w, r)

		Expect(w.Result().StatusCode).To(Equal(http.StatusCreated))

	})

	AfterEach(func() {
		filebasedDb.DB.Close()
	})

	Describe("Auth Middleware", func() {
		var (
			router *gin.Engine
			w      *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			router = gin.Default()
			w = httptest.NewRecorder()
		})

		When("valid token is provided", func() {
			It("should set user ID in context and call next middleware", func() {
				claims := &model.Claims{UserID: 123}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString(model.JwtKey)
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "session_token", Value: signedToken})

				router.Use(middleware.Auth())
				router.GET("/", func(ctx *gin.Context) {
					userID := ctx.MustGet("id").(int)
					Expect(userID).To(Equal(123))
				})

				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		When("session token is missing", func() {
			It("should return unauthorized error response", func() {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)

				router.Use(middleware.Auth())

				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusSeeOther))
			})
		})

		When("session token is invalid", func() {
			It("should return unauthorized error response", func() {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "session_token", Value: "invalid_token"})

				router.Use(middleware.Auth())

				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("Repository", func() {
		Describe("User", func() {
			When("fetching a single user data by email from users table in the database", func() {
				It("should return a single task data", func() {
					expectUser := model.User{
						Fullname: "test",
						Email:    "test@mail.com",
						Password: "testing123",
					}

					resUser, err := userRepo.GetUserByEmail("test@mail.com")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resUser.Fullname).To(Equal(expectUser.Fullname))
					Expect(resUser.Email).To(Equal(expectUser.Email))
					Expect(resUser.Password).To(Equal(expectUser.Password))
				})
			})

			When("retrieving user task categories from user repository", func() {
				It("should return the expected user task categories", func() {
					resUserTask, err := userRepo.GetUserTaskCategory()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resUserTask).To(Equal(expectedUserTask))
				})
			})
		})

		Describe("Task", func() {
			When("updating a category with valid ID and new category data", func() {
				It("should update the existing category with the new category data", func() {
					newCategory := model.Category{
						ID:   1,
						Name: "Updated with Repository Category 1",
					}
					err = categoryRepo.Update(1, newCategory)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := categoryRepo.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.Name).To(Equal(newCategory.Name))

					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("deleting a category with a valid category ID", func() {
				It("should delete the category from the database without returning an error", func() {
					err = categoryRepo.Delete(2)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := categoryRepo.GetByID(2)
					Expect(err.Error()).To(Equal("record not found"))
					Expect(result).To(BeNil())
				})
			})

			When("retrieving the list of categories from the database", func() {
				It("should return the list of categories without any errors and the list should contain the expected number of categories", func() {
					results, err := categoryRepo.GetList()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).To(HaveLen(5))

					Expect(results).To(Equal(insertCategories))
				})
			})

		})

		Describe("Task", func() {
			When("updating task data in tasks table in the database", func() {
				It("should update the existing task data in tasks table in the database", func() {
					newTask := model.Task{
						ID:         1,
						Title:      "Updated with Repository Task 1",
						Deadline:   "2023-05-30",
						Priority:   2,
						CategoryID: 1,
						Status:     "In Progress",
					}
					err = taskRepo.Update(&newTask)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := taskRepo.GetByID(1)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(result.Title).To(Equal(newTask.Title))
					Expect(result.Deadline).To(Equal(newTask.Deadline))
					Expect(result.Priority).To(Equal(newTask.Priority))
					Expect(result.CategoryID).To(Equal(newTask.CategoryID))
					Expect(result.Status).To(Equal(newTask.Status))
				})
			})

			When("deleting a task with a valid task ID from the database", func() {
				It("should delete the task without any errors", func() {
					err = taskRepo.Delete(2)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := taskRepo.GetByID(2)
					Expect(err.Error()).To(Equal("record not found"))
					Expect(result).To(BeNil())
				})
			})

			When("retrieving the list of tasks from the database", func() {
				It("should return the list of tasks without any errors", func() {
					results, err := taskRepo.GetList()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results).To(HaveLen(5))

					Expect(results).To(Equal(insertTasks))
				})
			})

			When("retrieving the list of tasks for a specific category from the database", func() {
				It("should return the list of tasks for the specified category without any errors", func() {
					taskCategory, err := taskRepo.GetTaskCategory(1)
					Expect(err).ShouldNot(HaveOccurred())

					Expect(taskCategory).To(Equal([]model.TaskCategory{
						{ID: 1, Title: "Task 1", Category: "Category 1"},
						{ID: 3, Title: "Task 3", Category: "Category 1"},
						{ID: 4, Title: "Task 4", Category: "Category 1"},
					}))
				})
			})

		})
	})

	Describe("Service", func() {
		Describe("User Service", func() {
			Describe("Login", func() {
				When("", func() {
					It("", func() {
						user := model.User{
							Fullname: "test",
							Email:    "test@mail.com",
							Password: "testing123",
						}
						_, err := userService.Login(&user)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
			Describe("GetUserTaskCategory", func() {
				When("retrieving user task categories from user repository", func() {
					It("should return the expected user task categories", func() {
						resUserTask, err := userService.GetUserTaskCategory()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resUserTask).To(Equal(expectedUserTask))
					})
				})
			})
		})

		Describe("Category Service", func() {
			Describe("Update", func() {
				When("updating a category in the database", func() {
					It("should update the category without any errors", func() {
						category := model.Category{
							ID:   1,
							Name: "Updated with Service Category 1",
						}

						err := categoryService.Update(1, category)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Delete", func() {
				When("deleting a category from the database", func() {
					It("should delete the category without any errors", func() {
						err := categoryService.Delete(3)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("GetList", func() {
				When("retrieving the list of categories from the database", func() {
					It("should return the list of categories without any errors", func() {
						categories, err := categoryService.GetList()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(categories).To(HaveLen(5))

						Expect(categories).To(Equal(insertCategories))
					})
				})
			})
		})

		Describe("Task Service", func() {
			Describe("Update", func() {
				When("updating a task in the database", func() {
					It("should update the task without any errors", func() {
						task := &model.Task{
							ID:         1,
							Title:      "Updated with Service Task 1",
							Deadline:   "2023-05-30",
							Priority:   5,
							CategoryID: 1,
							Status:     "In Progress",
						}

						err := taskService.Update(task)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("Delete", func() {
				When("deleting a task from the database", func() {
					It("should delete the task without any errors", func() {
						err := taskService.Delete(3)
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})

			Describe("GetList", func() {
				When("retrieving the list of tasks from the database", func() {
					It("should return the list of tasks without any errors", func() {
						tasks, err := taskService.GetList()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(tasks).To(Equal(insertTasks))
					})
				})
			})

			Describe("GetTaskCategory", func() {
				When("retrieving the category of a task from the database", func() {
					It("should return the task category without any errors", func() {
						taskCategories, err := taskService.GetTaskCategory(1)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(taskCategories).To(Equal([]model.TaskCategory{
							{ID: 1, Title: "Task 1", Category: "Category 1"},
							{ID: 3, Title: "Task 3", Category: "Category 1"},
							{ID: 4, Title: "Task 4", Category: "Category 1"},
						}))
					})
				})
			})
		})
	})

	Describe("API", func() {
		Describe("User API", func() {
			When("send empty email and password with POST method", func() {
				It("should return a bad request", func() {
					loginData := model.UserLogin{
						Email:    "",
						Password: "",
					}

					body, _ := json.Marshal(loginData)
					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/user/login", bytes.NewReader(body))
					r.Header.Set("Content-Type", "application/json")

					apiServer.ServeHTTP(w, r)

					errResp := model.ErrorResponse{}
					err := json.Unmarshal(w.Body.Bytes(), &errResp)
					Expect(err).To(BeNil())
					Expect(w.Result().StatusCode).To(Equal(http.StatusBadRequest))
					Expect(errResp.Error).To(Equal("invalid decode json"))
				})
			})

			When("send email and password with POST method", func() {
				It("should return a success", func() {
					loginData := model.UserLogin{
						Email:    "test@mail.com",
						Password: "testing123",
					}
					body, _ := json.Marshal(loginData)
					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/user/login", bytes.NewReader(body))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var resp = map[string]interface{}{}
					err = json.Unmarshal(w.Body.Bytes(), &resp)
					Expect(err).To(BeNil())
					Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
					Expect(resp["message"]).To(Equal("login success"))
				})
			})

			Describe("GetUserTaskCategory", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						r, _ := http.NewRequest("GET", "/user/tasks", nil)
						w := httptest.NewRecorder()
						r.Header.Set("Content-Type", "application/json")
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("retrieving user list by task and category", func() {
					It("should return status code 200 and task list", func() {
						r, _ := http.NewRequest("GET", "/user/tasks", nil)
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var userTasks []model.UserTaskCategory
						Expect(json.Unmarshal(w.Body.Bytes(), &userTasks)).Should(Succeed())
						Expect(userTasks).To(Equal(expectedUserTask))
					})
				})
			})
		})

		Describe("Category API", func() {
			Describe("UpdateCategory", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						updatedCategory := model.Category{
							ID:   1,
							Name: "Updated with API Category 1",
						}

						requestBody, _ := json.Marshal(updatedCategory)
						r, _ := http.NewRequest("PUT", "/category/update/1", bytes.NewReader(requestBody))
						r.Header.Set("Content-Type", "application/json")
						w := httptest.NewRecorder()
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("updating an existing category", func() {
					It("should update the category and return status code 200", func() {
						updatedCategory := model.Category{
							ID:   1,
							Name: "Updated with API Category 1",
						}

						requestBody, _ := json.Marshal(updatedCategory)
						r, _ := http.NewRequest("PUT", "/category/update/1", bytes.NewReader(requestBody))
						r.Header.Set("Content-Type", "application/json")
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response model.SuccessResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Message).To(Equal("category update success"))
					})
				})

				When("updating a non-existing category", func() {
					It("should return status code 400", func() {
						updatedCategory := model.Category{
							ID:   1,
							Name: "Updated with API Category 1",
						}

						requestBody, _ := json.Marshal(updatedCategory)
						r, _ := http.NewRequest("PUT", "/category/update/abc", bytes.NewReader(requestBody))
						r.Header.Set("Content-Type", "application/json")
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusBadRequest))

						var response model.ErrorResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Error).To(Equal("invalid Category ID"))
					})
				})

				When("sending an invalid request", func() {
					It("should return status code 400", func() {
						reqBody := []byte("invalid request body")

						r, _ := http.NewRequest("PUT", "/category/update/1", bytes.NewReader(reqBody))
						r.Header.Set("Content-Type", "application/json")
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusBadRequest))

						var response model.ErrorResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Error).NotTo(BeNil())
					})
				})
			})

			Describe("DeleteCategory", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						r, _ := http.NewRequest("DELETE", "/category/delete/4", nil)
						w := httptest.NewRecorder()
						r.Header.Set("Content-Type", "application/json")
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("deleting a category", func() {
					It("should delete the category and return status code 200", func() {
						r, _ := http.NewRequest("DELETE", "/category/delete/4", nil)
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response model.SuccessResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Message).To(Equal("category delete success"))
					})
				})
			})

			Describe("GetCategoryList", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						r, _ := http.NewRequest("GET", "/category/list", nil)
						w := httptest.NewRecorder()
						r.Header.Set("Content-Type", "application/json")
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("retrieving the list of categories", func() {
					It("should return the list of categories and status code 200", func() {
						r, _ := http.NewRequest("GET", "/category/list", nil)
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response []model.Category
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response).To(Equal(insertCategories))
					})
				})
			})
		})

		Describe("Task API", func() {
			Describe("UpdateTask", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						updatedTask := model.Task{
							ID:         1,
							Title:      "Updated with API Task 1",
							Deadline:   "2023-05-30",
							Priority:   5,
							CategoryID: 1,
							Status:     "In Progress",
						}
						reqBody, _ := json.Marshal(updatedTask)

						r, _ := http.NewRequest("PUT", fmt.Sprintf("/task/update/%d", 1), bytes.NewReader(reqBody))
						w := httptest.NewRecorder()
						r.Header.Set("Content-Type", "application/json")
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("updating existing task", func() {
					It("should return status code 200", func() {
						updatedTask := model.Task{
							ID:         1,
							Title:      "Updated with API Task 1",
							Deadline:   "2023-05-30",
							Priority:   5,
							CategoryID: 1,
							Status:     "In Progress",
						}
						reqBody, _ := json.Marshal(updatedTask)

						r, _ := http.NewRequest("PUT", fmt.Sprintf("/task/update/%d", 1), bytes.NewReader(reqBody))
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response model.SuccessResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Message).To(Equal("update task success"))
					})
				})

				When("sending invalid request", func() {
					It("should return status code 400", func() {
						reqBody := []byte("invalid request body")
						r, _ := http.NewRequest("PUT", "/task/update/1", bytes.NewReader(reqBody))
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)

						Expect(w.Code).To(Equal(http.StatusBadRequest))

						var response model.ErrorResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Error).NotTo(BeNil())
					})
				})
			})

			Describe("DeleteTask", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						r, _ := http.NewRequest("DELETE", fmt.Sprintf("/task/delete/%d", 1), nil)
						w := httptest.NewRecorder()
						r.Header.Set("Content-Type", "application/json")
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("deleting existing task", func() {
					It("should return status code 200", func() {
						r, _ := http.NewRequest("DELETE", fmt.Sprintf("/task/delete/%d", 1), nil)
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response model.SuccessResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Message).To(Equal("delete task success"))
					})
				})

				When("deleting non-existing task", func() {
					It("should return status code 400", func() {
						taskID := "abc"
						r, _ := http.NewRequest("DELETE", fmt.Sprintf("/task/delete/%s", taskID), nil)
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusBadRequest))

						var response model.ErrorResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Error).To(Equal("Invalid task ID"))
					})
				})
			})

			Describe("GetTaskList", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						r, _ := http.NewRequest("GET", "/task/list", nil)
						w := httptest.NewRecorder()
						r.Header.Set("Content-Type", "application/json")
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("retrieving task list", func() {
					It("should return status code 200 and task list", func() {
						r, _ := http.NewRequest("GET", "/task/list", nil)
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response []model.Task
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response).To(Equal(insertTasks))
					})
				})
			})

			Describe("GetTaskListByCategory", func() {
				When("sending without cookie", func() {
					It("should return status code 401", func() {
						r, _ := http.NewRequest("GET", fmt.Sprintf("/task/category/%d", 1), nil)
						w := httptest.NewRecorder()
						r.Header.Set("Content-Type", "application/json")
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusUnauthorized))
					})
				})

				When("retrieving task list by category", func() {
					It("should return status code 200 and task list", func() {
						r, _ := http.NewRequest("GET", fmt.Sprintf("/task/category/%d", 1), nil)
						w := httptest.NewRecorder()

						r.AddCookie(SetCookie(apiServer))
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response []model.TaskCategory
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response).To(Equal([]model.TaskCategory{
							{ID: 1, Title: "Task 1", Category: "Category 1"},
							{ID: 3, Title: "Task 3", Category: "Category 1"},
							{ID: 4, Title: "Task 4", Category: "Category 1"},
						}))
					})
				})
			})
		})
	})
})
