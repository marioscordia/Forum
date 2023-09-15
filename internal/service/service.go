package service

import (
	"net/http"
	"newforum/internal/form"
	"newforum/internal/store"
)

type Admin interface{
	GetUsers() ([]*store.User, error)
	GetRequests() ([]*store.User, error)
	AdminApproval(form *form.AdminApproval) error
	SnippetReports() ([]*store.Snippet, error)
	CommentReports() ([]*store.Comment, error)
	UnreportSnippet(snippetID int) error
	UnreportComment(commentID int) error
}

type Moderator interface{
	PendingSnippets() ([]*store.Snippet, error)
	Approval(form *form.Approval) error
	ReportSnippet(snippetID int) error
	ReportComment(form *form.Report) error
}

type User interface{
	CreateUser(form *form.UserSignup) error
	Authenticate(form *form.UserSignin) (int, error)
	CreateSession(w http.ResponseWriter, userID int) error
	DeleteSession(w http.ResponseWriter, cookie *http.Cookie) error
	TakeInfo(token string) (store.User, error)
	MakeRequest(userID int) error
}

type Comment interface {
	GetComments(snippetID int) ([]*store.Comment, error)
	GetComment(commentID int) (*store.Comment, error)
	MyComments(userID int) ([]*store.Comment, error)
	ReactedComments(userID int) ([]*store.Comment, error)
	CreateComment(form *form.Comment) (id int, err error)
	EditComment(form *form.CommentEdit) error
	DeleteComment(commentID int) error
	CommentReaction(form form.CommentReaction) error
}

type Snippet interface{
	GetSnippets() ([]*store.Snippet, error)
	FilterSnippets(filter form.Filter) ([]*store.Snippet, error)
	CreateSnippet(form *form.SnippetCreate) (int, error)
	EditSnippet(form *form.SnippetEdit) error
	DeleteSnippet(snippetID int, image string) error
	GetSnippet(snippetID int) (*store.Snippet, error)
	ReactedSnippets(userID int) ([]*store.Snippet, error)
	GetCreated(userID int) ([]*store.Snippet, error)
	SnippetReaction(form form.SnippetReaction) error
}

type Notification interface{
	Notifications(userID int) ([]*store.Notification, error)
	NotificationNum(userID int) (int, error)
	Update(userID int) error
}

type Service struct {
	Admin
	Moderator
	User
	Snippet
	Comment
	Notification
}

func NewService(store *store.Store) *Service {
	return &Service{
		NewAdminService(store),
		NewModeratorService(store),
		NewUserService(store),
		NewSnippetService(store),
		NewCommentService(store),
		NewNotificationService(store),
	}
}