package todo

type User struct {
	Id       string `json:"-" db:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type List struct {
	Id    string `json:"id" db:"id"`
	Title string `json:"title" db:"title" binding:"required"`
}

type UsersList struct {
	Id     string `db:"id"`
	UserId string `db:"user_id"`
	ListId string `db:"list_id"`
}

type Task struct {
	Id    string `json:"id" db:"id"`
	Title string `json:"title" db:"title" binding:"required"`
	Done  bool   `json:"done" db:"done"`
}

type ListsTasks struct {
	Id     string `db:"id"`
	ListId string `db:"list_id"`
	TaskId string `db:"task_id"`
}
