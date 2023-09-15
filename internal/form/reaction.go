package form

type SnippetReaction struct {
	UserID int
	ReactorName string
	SnippetID int
	Reaction string
}

type CommentReaction struct {
	UserID int
	CommentID int
	Reaction string
}
