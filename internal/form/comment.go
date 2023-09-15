package form

import "newforum/internal/validator"

type Comment struct {
	UserID int
	CommentatorName string 
	SnippetID int
	Comment string
	validator.Validator
}

type CommentEdit struct {
	CommentID int
	Comment string
	validator.Validator
}