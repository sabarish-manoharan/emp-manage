package models

import (
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model // gorm.Model includes fields like ID , Created At , Updated At , Deleted At automatically
	Name       string
	Age        int
	Role       string
	Email      string
	Id         int
}

type User struct{
	gorm.Model
	Name  string  `json:"name"`
	Email string `gorm:"unique" json:"email"`
	Password string  `json:"password"`
}
