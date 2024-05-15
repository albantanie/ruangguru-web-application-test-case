package model

type Category struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type Task struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	Title      string `json:"title"`
	Deadline   string `json:"deadline"`
	Priority   int    `json:"priority"`
	CategoryID int    `json:"category_id"`
	Status     string `json:"status"`
}

type TaskCategory struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}
