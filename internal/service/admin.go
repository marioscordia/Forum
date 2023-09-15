package service

import (
	"newforum/internal/form"
	"newforum/internal/store"
)

type AdminService struct {
	user store.UserStore
	snippet store.SnippetStore
	comment store.CommentStore
}

func NewAdminService(store *store.Store) *AdminService {
	return &AdminService{
		user: store.UserStore,
		snippet: store.SnippetStore,
		comment: store.CommentStore,
	}
}

func (s *AdminService) GetUsers() ([]*store.User, error){
	return s.user.GetUsers()
}

func (s *AdminService) GetRequests() ([]*store.User, error) {
	return s.user.GetRequests()
}

func (s *AdminService) SnippetReports() ([]*store.Snippet, error) {
	return s.snippet.SnippetReports()
}

func (s *AdminService) CommentReports() ([]*store.Comment, error) {
	return s.comment.CommentReports()
}

func (s *AdminService) UnreportSnippet(snippetID int) error{
	return s.snippet.UnreportSnippet(snippetID)
}

func (s *AdminService) UnreportComment(commentID int) error{
	return s.comment.UnreportComment(commentID)
}

func (s *AdminService) AdminApproval(form *form.AdminApproval) error{
	if form.Decision == "upgrade"{
		return s.user.Upgrade(form.UserID)
	}else if form.Decision == "downgrade"{
		return s.user.Downgrade(form.UserID)
	}
	
	return s.user.Reject(form.UserID)
}

