package models

type User struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Role  string `bson:"role"`
}
