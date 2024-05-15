package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
)

type TaskRepository interface {
	Store(task *model.Task) error
	Update(task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList() ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
}

type taskRepository struct {
	filebased *filebased.Data
}

func NewTaskRepo(filebasedDb *filebased.Data) *taskRepository {
	return &taskRepository{
		filebased: filebasedDb,
	}
}

func (t *taskRepository) Store(task *model.Task) error {
	t.filebased.StoreTask(*task)

	return nil
}

func (t *taskRepository) Update(task *model.Task) error {

	err := t.filebased.UpdateTask(task.ID, *task)

	if err != nil {
		return err
	}

	return nil // TODO: replace this
}

func (t *taskRepository) Delete(id int) error {

	err := t.filebased.DeleteTask(id)

	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	task, err := t.filebased.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	return task, nil // TODO: replace this
}

func (t *taskRepository) GetList() ([]model.Task, error) {
	task, err := t.filebased.GetTasks()
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	Category, err := t.filebased.GetTaskListByCategory(id)

	if err != nil {
		return nil, err
	}

	return Category, nil
}
