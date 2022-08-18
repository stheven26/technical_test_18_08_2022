package models

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username" binding:"required" gorm:"unique"`
	Password string `json:"-" binding:"required"`
}

type Blog struct {
	Id    uint   `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
	Slug  string `json:"slug" binding:"required"`
}
