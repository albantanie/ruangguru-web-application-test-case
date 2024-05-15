package main_test

import (
	main "a21hc3NpZ25tZW50"
	"a21hc3NpZ25tZW50/db/filebased"
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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task Tracker Plus", Ordered, func() {
	var apiServer *gin.Engine

	var categoryRepo repo.CategoryRepository
	var taskRepo repo.TaskRepository

	var categoryService service.CategoryService
	var taskService service.TaskService

	var insertCategories []model.Category
	var insertTasks []model.Task

	filebasedDb, err := filebased.InitDB()

	Expect(err).ShouldNot(HaveOccurred())

	categoryRepo = repo.NewCategoryRepo(filebasedDb)
	taskRepo = repo.NewTaskRepo(filebasedDb)

	categoryService = service.NewCategoryService(categoryRepo)
	taskService = service.NewTaskService(taskRepo)

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode) //release

		os.Remove("file.db")

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
				CategoryID: 1,
				Status:     "In Progress",
			},
			{
				ID:         2,
				Title:      "Task 2",
				Deadline:   "2023-06-01",
				Priority:   1,
				CategoryID: 2,
				Status:     "Completed",
			},
			{
				ID:         3,
				Title:      "Task 3",
				Deadline:   "2023-06-02",
				Priority:   4,
				CategoryID: 1,
				Status:     "Completed",
			},
			{
				ID:         4,
				Title:      "Task 4",
				Deadline:   "2023-06-02",
				Priority:   3,
				CategoryID: 1,
				Status:     "Completed",
			},
			{
				ID:         5,
				Title:      "Task 5",
				Deadline:   "2023-06-07",
				Priority:   5,
				CategoryID: 3,
				Status:     "In Progress",
			},
		}

		for _, v := range insertTasks {
			err := taskRepo.Store(&v)
			Expect(err).ShouldNot(HaveOccurred())
		}

	})

	Describe("Repository", func() {
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

					Expect(err).ShouldNot(HaveOccurred())
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

		BeforeEach(func() {
			apiServer = gin.New()
			apiServer = main.RunServer(apiServer, filebasedDb)
		})
		Describe("Category API", func() {
			Describe("UpdateCategory", func() {
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
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusBadRequest))

						var response model.ErrorResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Error).NotTo(BeNil())
					})
				})
			})

			Describe("DeleteCategory", func() {
				When("deleting a category", func() {
					It("should delete the category and return status code 200", func() {
						r, _ := http.NewRequest("DELETE", "/category/delete/4", nil)
						w := httptest.NewRecorder()
						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response model.SuccessResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Message).To(Equal("category delete success"))
					})
				})
			})

			Describe("GetCategoryList", func() {
				When("retrieving the list of categories", func() {
					It("should return the list of categories and status code 200", func() {
						r, _ := http.NewRequest("GET", "/category/list", nil)
						w := httptest.NewRecorder()

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
			Describe("/task/update/:id", func() {
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

						apiServer.ServeHTTP(w, r)

						Expect(w.Code).To(Equal(http.StatusBadRequest))

						var response model.ErrorResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Error).NotTo(BeNil())
					})
				})
			})

			Describe("/task/delete/:id", func() {
				When("deleting existing task", func() {
					It("should return status code 200", func() {
						r, _ := http.NewRequest("DELETE", fmt.Sprintf("/task/delete/%d", 1), nil)
						w := httptest.NewRecorder()

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

						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusBadRequest))

						var response model.ErrorResponse
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response.Error).To(Equal("Invalid task ID"))
					})
				})
			})

			Describe("/task/list", func() {
				When("retrieving task list", func() {
					It("should return status code 200 and task list", func() {
						r, _ := http.NewRequest("GET", "/task/list", nil)
						w := httptest.NewRecorder()

						apiServer.ServeHTTP(w, r)
						Expect(w.Code).To(Equal(http.StatusOK))

						var response []model.Task
						Expect(json.Unmarshal(w.Body.Bytes(), &response)).Should(Succeed())
						Expect(response).To(Equal(insertTasks))
					})
				})
			})

			Describe("/task/category/:id", func() {
				When("retrieving task list by category", func() {
					It("should return status code 200 and task list", func() {
						r, _ := http.NewRequest("GET", fmt.Sprintf("/task/category/%d", 1), nil)
						w := httptest.NewRecorder()

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
