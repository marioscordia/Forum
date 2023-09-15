package form

import "newforum/internal/validator"

type UserSignup struct{
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
	Confirm string `form:"confirm"`
	validator.Validator
}

type UserSignin struct{
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator
}

type UserInfo struct{
	ID int
	Name string
}

type AdminApproval struct{
	UserID int
	Decision string
}