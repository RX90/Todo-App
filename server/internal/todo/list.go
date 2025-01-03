package todo

type List struct {
	Id    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title" binding:"required"`
}

type UsersList struct {
	Id     int `db:"id"`
	UserId int `db:"user_id"`
	ListId int `db:"list_id"`
}
