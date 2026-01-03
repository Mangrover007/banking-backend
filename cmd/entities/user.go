package entities

type User struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	ID       string `json:"-" binding:"-"`
}
