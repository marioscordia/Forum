package store

import (
	"database/sql"
	"newforum/internal/form"
)

type UserStore interface{
	GetUsers() ([]*User, error)
	Insert(form *form.UserSignup) error
	Authenticate(form *form.UserSignin) (int, error)
	PutToken(userID int, token string) error
	TakeInfo(token string) (User, error)
	DeleteToken(token string) error
	UserExists(email string) (bool, error)
	MakeRequest(userID int) error
	GetRequests() ([]*User, error)
	Reject(userID int) error
	Upgrade(userID int) error
	Downgrade(userID int) error
}

type SnippetStore interface{
	LatestSnippets() ([]*Snippet, error)
	PendingSnippets() ([]*Snippet, error)
	SnippetReports() ([]*Snippet, error)
	FilterSnippets(category form.Filter) ([]*Snippet, error)
	CreatedSnippets(userID int) ([]*Snippet, error)
	ReactedSnippets(userID int) ([]*Snippet, error)
	CreateSnippet(form *form.SnippetCreate) (int, error)
	UpdateSnippet(form *form.SnippetEdit) error
	GetSnippet(snippetID int) (*Snippet, error)
	ApproveSnippet(snippetID int) error
	ReportSnippet(snippetID int) error
	UnreportSnippet(snippetID int) error
	DeleteSnippet(snippetID int) error
}

type ReactionStore interface{
	PostReaction(form form.SnippetReaction) error
	ExistsReactionPost(snippetID, userID int, reaction string) (bool, error)
	DeleteReactionPost(snippetID, userID int, reaction string) error
	CommentReaction(form form.CommentReaction) error
	ExistsReactionComment(commentID, userID int, reaction string) (bool, error)
	DeleteReactionComment(commentID, userID int, reaction string) error
}

type CommentStore interface{
	CreateComment(form *form.Comment) (int, error)
	GetComments(snippetID int) ([]*Comment, error)
	CommentReports() ([]*Comment, error)
	GetComment(commentID int) (*Comment, error)
	ReportComment(form *form.Report) error
	UnreportComment(commentID int) error
	EditComment(form *form.CommentEdit) error
	DeleteComment(commentID int) error
	ReactedComments(userID int) ([]*Comment, error)
	MyComments(userID int) ([]*Comment, error)
}

type NotificationStore interface{
	Notifications(userID int) ([]*Notification, error)
	Update(userID int) error
	NotificationNum(userID int) (int, error)
}

type Store struct{
	UserStore
	SnippetStore
	ReactionStore
	CommentStore
	NotificationStore
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		NewUserModel(db),
		NewSnippetModel(db),
		NewReactionModel(db),
		NewCommentModel(db),
		NewNotificationsModel(db),
	}
}