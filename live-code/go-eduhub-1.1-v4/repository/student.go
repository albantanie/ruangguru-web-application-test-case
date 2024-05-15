package repository

import "a21hc3NpZ25tZW50/model"

type StudentRepository interface {
	Store(student *model.Student) error
}

type studentRepository struct {
	students []model.Student
}

func NewStudentRepo() *studentRepository {
	return &studentRepository{}
}

func (s *studentRepository) Store(student *model.Student) error {
	s.students = append(s.students, *student)
	return nil // TODO: Replace this with appropriate error handling
}
