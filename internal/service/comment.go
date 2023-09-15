package service

import (
	"newforum/internal/form"
	"newforum/internal/oops"
	"newforum/internal/store"
	"newforum/internal/validator"
)

type CommentService struct {
	comment store.CommentStore
	reaction store.ReactionStore
}

func NewCommentService(store *store.Store) *CommentService {
	return &CommentService{
		comment: store.CommentStore,
		reaction: store.ReactionStore,
	}
}

func (s *CommentService) GetComments(id int) ([]*store.Comment, error) {
	return s.comment.GetComments(id)
}

func (s *CommentService) GetComment(commentID int) (*store.Comment, error) {
	return s.comment.GetComment(commentID)
}

func (s *CommentService) MyComments(userID int) ([]*store.Comment, error) {
	return s.comment.MyComments(userID)
}

func (s *CommentService) ReactedComments(userID int) ([]*store.Comment, error) {
	return s.comment.ReactedComments(userID)
}

func (s *CommentService) CreateComment(form *form.Comment) (int, error) {
	form.CheckField(validator.NotBlank(form.Comment), "comment", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Comment, 100), "comment", "This field cannot be more than 100 characters long")

	if !form.Valid(){
		return 0, oops.ErrFormInvalid
	}
	
	return s.comment.CreateComment(form)
}

func (s *CommentService) EditComment(form *form.CommentEdit) error {
	form.CheckField(validator.NotBlank(form.Comment), "comment", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Comment, 100), "comment", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		return oops.ErrFormInvalid
	}

	return s.comment.EditComment(form)
}

func (s *CommentService) DeleteComment(commentID int) error {
	return s.comment.DeleteComment(commentID)
}

func (s *CommentService) CommentReaction(form form.CommentReaction) error {
	return s.reaction.CommentReaction(form)
}