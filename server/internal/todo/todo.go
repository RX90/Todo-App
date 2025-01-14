package todo

import "errors"

type User struct {
	Id       string `json:"-"        db:"id"`
	Username string `json:"username" db:"username"      binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
}

type List struct {
	Id    string `json:"id"    db:"id"`
	Title string `json:"title" db:"title" binding:"required"`
}

type UsersList struct {
	Id     string
	UserId string
	ListId string
}

type Task struct {
	Id    string `json:"id"    db:"id"`
	Title string `json:"title" db:"title"`
	Done  bool   `json:"done"  db:"done"`
}

type ListsTasks struct {
	Id     string
	ListId string
	TaskId string
}

type UpdateTaskInput struct {
	Title *string `json:"title" db:"title"`
	Done  *bool   `json:"done"  db:"done"`
}

func (i UpdateTaskInput) Validate() error {
	if i.Title == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
