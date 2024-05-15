package main_test

import (
	main "a21hc3NpZ25tZW50"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go EduHub 1", Ordered, func() {
	var apiServer *gin.Engine
	var studentRepo repo.StudentRepository
	var courseRepo repo.CourseRepository

	studentRepo = repo.NewStudentRepo()
	courseRepo = repo.NewCourseRepo()

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode) //release

		apiServer = gin.New()
		apiServer = main.RunServer(apiServer)
	})

	Describe("Repository", func() {
		Describe("Student repository", func() {
			When("add student data to students inmemory database", func() {
				It("should save student data to students inmemory database", func() {
					student := model.Student{
						Name:     "John",
						Email:    "Jl. Raya Cianjur",
						Phone:    "08342424275",
						CourseID: 1,
					}
					err := studentRepo.Store(&student)
					Expect(err).ShouldNot(HaveOccurred())

					results, err := studentRepo.FetchAll()
					Expect(err).ShouldNot(HaveOccurred())

					Expect(results[0].Name).To(Equal(student.Name))
					Expect(results[0].Email).To(Equal(student.Email))
					Expect(results[0].Phone).To(Equal(student.Phone))
					Expect(results[0].CourseID).To(Equal(student.CourseID))

					studentRepo.ResetStudentRepo()
				})
			})

			When("fetching a single student data by id from students table in the database", func() {
				It("should return a single student data", func() {
					student := model.Student{
						Name:     "John",
						Email:    "Jl. Raya Cilandak",
						Phone:    "08121322432",
						CourseID: 1,
					}
					err := studentRepo.Store(&student)
					Expect(err).ShouldNot(HaveOccurred())

					results, err := studentRepo.FetchAll()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(results[0].Name).To(Equal(student.Name))
					Expect(results[0].Email).To(Equal(student.Email))
					Expect(results[0].Phone).To(Equal(student.Phone))
					Expect(results[0].CourseID).To(Equal(student.CourseID))

					studentRepo.ResetStudentRepo()
				})
			})
		})

		Describe("Course repository", func() {
			When("add course data to course inmemory database", func() {
				It("should save course data to courses inmemory database", func() {
					course := model.Course{
						ID:         1,
						Name:       "Introduction to Computer Science",
						Schedule:   "Monday and Wednesday, 10am - 12pm",
						Grade:      3.5, // (dalam skala 4.0)
						Attendance: 85,  // (dalam persen)
					}
					err := courseRepo.Store(&course)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := courseRepo.FetchByID(1)
					Expect(err).ShouldNot(HaveOccurred())

					Expect(result.Name).To(Equal(course.Name))
					Expect(result.Schedule).To(Equal(course.Schedule))
					Expect(result.Grade).To(Equal(course.Grade))
					Expect(result.Attendance).To(Equal(course.Attendance))

					courseRepo.ResetCourseRepo()
				})
			})
		})
	})

	Describe("API", func() {
		Describe("/student/add", func() {
			When("sending valid request", func() {
				It("should return status code 200", func() {
					student := model.Student{
						Name:     "Aditira Jamhuri",
						Email:    "aditira@gmail.com",
						Phone:    "085231314132",
						CourseID: 1,
					}

					body, _ := json.Marshal(student)
					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/student/add", bytes.NewReader(body))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var response model.SuccessResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Message).To(Equal("add student success"))

					studentRepo.ResetStudentRepo()
				})
			})

			When("sending invalid request", func() {
				It("should return status code 400", func() {
					reqBody := []byte("invalid request body")

					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/student/add", bytes.NewReader(reqBody))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var response model.ErrorResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusBadRequest))
					Expect(response.Error).NotTo(BeNil())
				})
			})

			When("encountering internal server error", func() {
				It("should return status code 500", func() {
					reqBody, _ := json.Marshal(model.Invalid{
						ID:      1,
						Message: "Invalid Add",
					})

					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/student/add", bytes.NewReader(reqBody))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var response model.ErrorResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Error).NotTo(BeNil())
				})
			})
		})

		Describe("/student/get/:id", func() {
			var expectedStudents []model.Student
			BeforeEach(func() {
				expectedStudents = []model.Student{
					{
						ID:       1,
						Name:     "Aditira Jamhuri",
						Email:    "aditira@gmail.com",
						Phone:    "085231314132",
						CourseID: 1,
					},
					{
						ID:       2,
						Name:     "Dito",
						Email:    "dito@gmail.com",
						Phone:    "081232332323",
						CourseID: 2,
					},
					{
						ID:       3,
						Name:     "Imam",
						Email:    "imam@gmail.com",
						Phone:    "083231000009",
						CourseID: 2,
					},
				}

				for _, v := range expectedStudents {
					body, _ := json.Marshal(v)
					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/student/add", bytes.NewReader(body))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)
					Expect(w.Code).To(Equal(http.StatusOK))
				}
			})

			When("given a valid student ID", func() {
				It("should return the corresponding student", func() {
					w := httptest.NewRecorder()
					r := httptest.NewRequest("GET", "/student/get/3", nil)
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					Expect(w.Code).To(Equal(http.StatusOK))
					var respStudent model.Student
					err := json.Unmarshal(w.Body.Bytes(), &respStudent)
					Expect(err).To(BeNil())
					Expect(respStudent).To(Equal(expectedStudents[2]))
				})
			})

			When("when given an invalid student ID", func() {
				It("should return an error response", func() {
					w := httptest.NewRecorder()
					r := httptest.NewRequest("GET", "/student/get/abc", nil)
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					Expect(w.Code).To(Equal(http.StatusBadRequest))
					var errResp model.ErrorResponse
					err := json.Unmarshal(w.Body.Bytes(), &errResp)
					Expect(err).To(BeNil())
					Expect(errResp.Error).To(Equal("Invalid student ID"))
				})
			})

			When("when the requested student does not exist", func() {
				It("should return a not found error response", func() {
					w := httptest.NewRecorder()
					r := httptest.NewRequest("GET", "/student/get/4", nil)
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					Expect(w.Code).To(Equal(http.StatusNotFound))
					var errResp model.ErrorResponse
					err := json.Unmarshal(w.Body.Bytes(), &errResp)
					Expect(err).To(BeNil())

					studentRepo.ResetStudentRepo()
				})
			})
		})

		Describe("/course/add", func() {
			When("sending valid request", func() {
				It("should return status code 200", func() {
					course := model.Course{
						Name:       "Introduction to Computer Science",
						Schedule:   "Monday and Wednesday, 10am - 12pm",
						Grade:      3.5, // (dalam skala 4.0)
						Attendance: 85,  // (dalam persen)
					}

					body, _ := json.Marshal(course)
					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/course/add", bytes.NewReader(body))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var response model.SuccessResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Message).To(Equal("add course success"))

					courseRepo.ResetCourseRepo()
				})
			})

			When("sending invalid request", func() {
				It("should return status code 400", func() {
					reqBody := []byte("invalid request body")

					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/course/add", bytes.NewReader(reqBody))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var response model.ErrorResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusBadRequest))
					Expect(response.Error).NotTo(BeNil())
				})
			})

			When("encountering internal server error", func() {
				It("should return status code 500", func() {
					reqBody, _ := json.Marshal(model.Invalid{
						ID:      1,
						Message: "Invalid Add",
					})

					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/course/add", bytes.NewReader(reqBody))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var response model.ErrorResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Error).NotTo(BeNil())
				})
			})
		})
	})
})
