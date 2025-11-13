package types

type User struct {
	Id int64
	Name string `validate:"required"`
	Email string `validate:"required"`
	Age int `validate:"required"`
}