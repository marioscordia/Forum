package service

import (
	"fmt"
	"io"
	"newforum/internal/form"
	"newforum/internal/oops"
	"newforum/internal/store"
	"newforum/internal/validator"
	"os"
	"time"
)

type SnippetService struct{
	snippet store.SnippetStore
	comment store.CommentStore
	reaction store.ReactionStore
}

func NewSnippetService(store *store.Store) *SnippetService {
	return &SnippetService{
		snippet: store.SnippetStore,
		comment: store.CommentStore,
		reaction: store.ReactionStore,
	}
}

func (s *SnippetService) GetSnippets() ([]*store.Snippet, error) {
	return s.snippet.LatestSnippets()
}

func (s *SnippetService) FilterSnippets(form *form.Filter) ([]*store.Snippet, error) {
	form.CheckField(validator.CheckCategory(form.Category), "filter", "Choose filter")
	if !form.Valid() {
		return nil, oops.ErrFormInvalid
	}
	
	return s.snippet.FilterSnippets(*form)
}

func (s *SnippetService) CreateSnippet(form *form.SnippetCreate) (int, error) {
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 50), "title", "This field cannot be more than 50 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 300), "content", "This field cannot be more than 300 characters long")
	form.CheckField(validator.CheckCategory(form.Category), "category", "Choose at least one category")
	form.CheckField(validator.CheckFileName(form.FileName), "image", "File type not supported. Only .jpg, .jpeg, .gif, .png files are allowed.")
	form.CheckField(validator.CheckFileSize(int(form.FileSize)), "image", "File size exceeds the limit of 20 MB")
	
	if !form.Valid() {
		return 0, oops.ErrFormInvalid
	}

	imageUrl := "./internal/store/img/"
	timestamp := time.Now().String()
	form.FileName = timestamp + form.FileName

	id, err := s.snippet.CreateSnippet(form)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(imageUrl+form.FileName, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	io.Copy(f, form.File)

	return id, nil
}

func (s *SnippetService) EditSnippet(form *form.SnippetEdit) error {
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 50), "title", "This field cannot be more than 50 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 300), "content", "This field cannot be more than 300 characters long")

	if !form.Valid() {
		return oops.ErrFormInvalid
	}

	return s.snippet.UpdateSnippet(form)
}

func (s *SnippetService) GetSnippet(id int) (*store.Snippet, error) {
	return s.snippet.GetSnippet(id)
}

func (s *SnippetService) ReactedSnippets(userID int) ([]*store.Snippet, error) {
	return s.snippet.ReactedSnippets(userID)
}

func (s *SnippetService) GetCreated(userID int) ([]*store.Snippet, error) {
	return s.snippet.CreatedSnippets(userID)
}

func (s *SnippetService) SnippetReaction(form form.SnippetReaction) error {
	return s.reaction.PostReaction(form)
}

func (s *SnippetService) DeleteSnippet(snippetID int, image string) error {
	err := os.Remove(fmt.Sprintf("./internal/store/img/%s", image))
	if err != nil {
		return err
	}

	return s.snippet.DeleteSnippet(snippetID)
}



