package form

import (
	"mime/multipart"
	"newforum/internal/validator"
)

type SnippetCreate struct {
	UserID int 
	Title string
	Content string
	Category []string
	FileName string
	FileSize int
	File multipart.File
	validator.Validator
}

type SnippetEdit struct { 
	SnippetID int
	Title string
	Content string
	validator.Validator 
}

type Approval struct {
	SnippetID int
	Image string
	Decision string
}

type Report struct {
	SnippetID int
	CommentID int
}