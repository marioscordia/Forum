package service

import (
	"fmt"
	"newforum/internal/form"
	"newforum/internal/store"
	"os"
)

type ModeratorService struct {
	snippet store.SnippetStore
	comment store.CommentStore
}

func NewModeratorService(store *store.Store) *ModeratorService {
	return &ModeratorService{
		snippet: store.SnippetStore,
		comment: store.CommentStore,
	}
}

func (s *ModeratorService) PendingSnippets() ([]*store.Snippet, error){
	return s.snippet.PendingSnippets()
}

func (s *ModeratorService) Approval(form *form.Approval) error{
	if form.Decision == "approve"{
		return s.snippet.ApproveSnippet(form.SnippetID)
	}
	err := os.Remove(fmt.Sprintf("./internal/store/img/%s", form.Image))
	if err != nil{
		return err
	}
	return s.snippet.DeleteSnippet(form.SnippetID)
}

func (s *ModeratorService) ReportSnippet(snippetID int) error {
	return s.snippet.ReportSnippet(snippetID)
}

func (s *ModeratorService) ReportComment(form *form.Report) error {
	return s.comment.ReportComment(form)
}