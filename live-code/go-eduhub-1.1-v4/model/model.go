package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	CourseID int    `json:"course_id"`
}

type Course struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Schedule   string  `json:"schedule"`
	Grade      float32 `json:"grade,omitempty"`
	Attendance int     `json:"attendance,omitempty"`
}

type Invalid struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}
