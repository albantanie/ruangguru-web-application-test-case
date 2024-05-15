package main_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	main "a21hc3NpZ25tZW50"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func setHeader(req *http.Request) {
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("user1:pass1")))
}

var _ = Describe("Blog Middleware", func() {
	var (
		resp   *httptest.ResponseRecorder
		router *gin.Engine
	)

	BeforeEach(func() {
		router = main.SetupRouter()
	})

	Describe("GET /posts is called", func() {
		When("user is authenticated", func() {
			It("returns HTTP status code 200 and a JSON array of posts", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
				resp = httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))
				var data map[string][]main.Post
				json.Unmarshal(resp.Body.Bytes(), &data)
				Expect(data["posts"]).To(HaveLen(len(main.Posts)))
				Expect(data["posts"][0].ID).To(Equal(main.Posts[0].ID))
				Expect(data["posts"][0].Title).To(Equal(main.Posts[0].Title))
				Expect(data["posts"][0].Content).To(Equal(main.Posts[0].Content))
			})
		})

		When("user is not authenticated", func() {
			It("should return 401 Unauthorized", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
				resp = httptest.NewRecorder()
				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("GET /posts?id={id} is called with a valid id", func() {
		When("user is authenticated", func() {
			It("returns HTTP status code 200 and a JSON object of the post with the given id", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts?id=1", nil)
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusOK))
				var data map[string]main.Post
				json.Unmarshal(resp.Body.Bytes(), &data)
				Expect(data["post"].ID).To(Equal(main.Posts[0].ID))
			})
		})

		When("user is not authenticated", func() {
			It("should return 401 Unauthorized", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts?id=1", nil)
				resp = httptest.NewRecorder()
				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("GET /posts?id={id} is called with an invalid id", func() {
		When("user is authenticated and number query param", func() {
			It("returns HTTP status code 404", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts?id=999", nil)
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusNotFound))
			})

			It("returns an error message in JSON format", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts?id=999", nil)
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				var data map[string]string
				json.Unmarshal(resp.Body.Bytes(), &data)
				Expect(data["error"]).To(Equal("Postingan tidak ditemukan"))
			})
		})

		When("user is authenticated and not number query param", func() {
			It("returns HTTP status code 400", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts?id=A", nil)
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("returns an error message in JSON format", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts?id=A", nil)
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				var data map[string]string
				json.Unmarshal(resp.Body.Bytes(), &data)
				Expect(data["error"]).To(Equal("ID harus berupa angka"))
			})
		})

		When("user is not authenticated", func() {
			It("should return 401 Unauthorized", func() {
				req, _ := http.NewRequest(http.MethodGet, "/posts?id=999", nil)
				resp = httptest.NewRecorder()
				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("POST /posts is called with valid JSON data", func() {
		When("user is authenticated", func() {
			It("returns HTTP status code 201", func() {
				post := main.Post{Title: "Judul Postingan Baru", Content: "Ini adalah postingan baru di blog ini."}
				body, _ := json.Marshal(post)
				req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusCreated))
			})

			It("adds the new post to the list of posts", func() {
				post := main.Post{Title: "Judul Postingan Baru", Content: "Ini adalah postingan baru di blog ini."}
				body, _ := json.Marshal(post)
				req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(len(main.Posts)).To(Equal(4))
				Expect(main.Posts[3].ID).To(Equal(4))
				Expect(main.Posts[3].Title).To(Equal("Judul Postingan Baru"))
				Expect(main.Posts[3].Content).To(Equal("Ini adalah postingan baru di blog ini."))
			})

			It("returns the new post in JSON format", func() {
				post := main.Post{Title: "Judul Postingan Baru", Content: "Ini adalah postingan baru di blog ini."}
				body, _ := json.Marshal(post)
				req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				var data map[string]main.Post
				json.Unmarshal(resp.Body.Bytes(), &data)
				Expect(data["post"].Title).To(Equal("Judul Postingan Baru"))
				Expect(data["post"].Content).To(Equal("Ini adalah postingan baru di blog ini."))
			})
		})

		When("user is not authenticated", func() {
			It("should return 401 Unauthorized", func() {
				req, _ := http.NewRequest(http.MethodPost, "/posts", nil)
				resp = httptest.NewRecorder()
				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})

	Describe("POST /posts is called with invalid JSON datad", func() {
		When("user is authenticated", func() {
			It("returns HTTP status code 400", func() {
				body := []byte(`{"title": "Judul Postingan Baru", "content":}`)
				req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("does not add a new post to the list of posts", func() {
				body := []byte(`{"title": "Judul Postingan Baru", "content":}`)
				req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				Expect(len(main.Posts)).To(Equal(5))
			})

			It("returns an error message in JSON format", func() {
				body := []byte(`{"title": "Judul Postingan Baru", "content":}`)
				req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				resp := httptest.NewRecorder()

				setHeader(req)

				router.ServeHTTP(resp, req)
				var data map[string]string
				json.Unmarshal(resp.Body.Bytes(), &data)
				Expect(data["error"]).To(Equal("Invalid request body"))
			})
		})

		When("user is not authenticated", func() {
			It("should return 401 Unauthorized", func() {
				req, _ := http.NewRequest(http.MethodPost, "/posts", nil)
				resp = httptest.NewRecorder()
				router.ServeHTTP(resp, req)
				Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})
})
