package filebased

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"a21hc3NpZ25tZW50/model"

	"go.etcd.io/bbolt"
)

type Data struct {
	DB *bbolt.DB
}

func InitDB() (*Data, error) {
	db, err := bbolt.Open("file.db", 0600, &bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			return fmt.Errorf("create tasks bucket: %v", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Categories"))
		if err != nil {
			return fmt.Errorf("create categories bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Data{DB: db}, nil
}

func (data *Data) StoreTask(task model.Task) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		return b.Put([]byte(fmt.Sprintf("%d", task.ID)), taskJSON)
	})
}

func (data *Data) StoreCategory(category model.Category) error {
	categoryJSON, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		return b.Put([]byte(fmt.Sprintf("%d", category.ID)), categoryJSON)
	})
}

func (data *Data) UpdateTask(id int, task model.Task) error {
	return data.StoreTask(task) // Reuse StoreTask as it will replace the existing entry
}

func (data *Data) UpdateCategory(id int, category model.Category) error {
	return data.StoreCategory(category) // Reuse StoreCategory as it will replace the existing entry
}

func (data *Data) DeleteTask(id int) error {
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		return b.Delete([]byte(fmt.Sprintf("%d", id)))
	})
}

func (data *Data) DeleteCategory(id int) error {
	return data.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		return b.Delete([]byte(fmt.Sprintf("%d", id)))
	})
}

func (data *Data) GetTaskByID(id int) (*model.Task, error) {
	var task model.Task
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		v := b.Get([]byte(fmt.Sprintf("%d", id)))
		if v == nil {
			return fmt.Errorf("record not found")
		}
		return json.Unmarshal(v, &task)
	})
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (data *Data) GetCategoryByID(id int) (*model.Category, error) {
	var category model.Category
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		v := b.Get([]byte(fmt.Sprintf("%d", id)))
		if v == nil {
			return fmt.Errorf("record not found")
		}
		return json.Unmarshal(v, &category)
	})
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (data *Data) GetTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		return b.ForEach(func(k, v []byte) error {
			var task model.Task
			if err := json.Unmarshal(v, &task); err != nil {
				log.Println("Error unmarshaling task:", err)
				return nil // Continue despite error
			}
			tasks = append(tasks, task)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %v", err)
	}
	return tasks, nil
}

func (data *Data) GetCategories() ([]model.Category, error) {
	var categories []model.Category
	err := data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Categories"))
		return b.ForEach(func(k, v []byte) error {
			var category model.Category
			if err := json.Unmarshal(v, &category); err != nil {
				log.Println("Error unmarshaling category:", err)
				return nil // Continue despite error
			}
			categories = append(categories, category)
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %v", err)
	}
	return categories, nil
}

func (data *Data) Reset() error {
	return data.DB.Update(func(tx *bbolt.Tx) error {
		if err := tx.DeleteBucket([]byte("Tasks")); err != nil {
			return err
		}
		if err := tx.DeleteBucket([]byte("Categories")); err != nil {
			return err
		}
		_, err := tx.CreateBucket([]byte("Tasks"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket([]byte("Categories"))
		return err
	})
}

func (data *Data) CloseDB() error {
	return data.DB.Close()
}

func (data *Data) GetTaskListByCategory(categoryID int) ([]model.TaskCategory, error) {
	var taskCategories []model.TaskCategory
	category, err := data.GetCategoryByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("error fetching category: %v", err)
	}

	err = data.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Tasks"))
		if b == nil {
			return fmt.Errorf("tasks bucket not found")
		}
		return b.ForEach(func(k, v []byte) error {
			var task model.Task
			if err := json.Unmarshal(v, &task); err != nil {
				log.Printf("Error unmarshaling task: %v", err)
				return nil // Continue processing next item in case of error
			}
			if task.CategoryID == categoryID {
				taskCategories = append(taskCategories, model.TaskCategory{
					ID:       task.ID,
					Title:    task.Title,
					Category: category.Name,
				})
			}
			return nil
		})
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks for category %d: %v", categoryID, err)
	}
	if len(taskCategories) == 0 {
		return nil, fmt.Errorf("no tasks found for category ID: %d", categoryID)
	}
	return taskCategories, nil
}
