package repository

import "a21hc3NpZ25tZW50/model"

type CourseRepository interface {
	Store(course *model.Course) error
}

type courseRepository struct {
	courses []model.Course
}

func NewCourseRepo() *courseRepository {
	return &courseRepository{}
}

func (c *courseRepository) Store(course *model.Course) error {
	c.courses = append(c.courses, *course)
	return nil // TODO: Replace this with appropriate error handling
}
