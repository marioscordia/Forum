package store

import (
	"database/sql"
	"errors"
	"newforum/internal/form"
	"time"
)

type Comment struct {
	ID int
	SnippetID int
	SnippetTitle string
	Content string
	Created time.Time
	Author string
	AuthorID int
	Likes int
	Dislikes int
	Status string
	Reported int
}

type CommentModel struct {
	DB *sql.DB
}

func NewCommentModel(db *sql.DB) *CommentModel {
	return &CommentModel{DB: db}
}

func (m *CommentModel) CreateComment(form *form.Comment) (int, error){
	stmt := `
	INSERT INTO comments (snippet_id, user_id, content, created)
	VALUES (?, ?, ?, datetime('now','+6 hours'));
	`
	res, err := m.DB.Exec(stmt, form.SnippetID, form.UserID, form.Comment)
	if err != nil {
		return 0, err
	}

	commentID, err := res.LastInsertId()
	if err != nil{
		return 0, err
	}

	var receiverID int

	stmt = `
	SELECT user_id FROM snippets
	WHERE id = ?;
	`
	err = m.DB.QueryRow(stmt, form.SnippetID).Scan(&receiverID)
	if err != nil {
		return 0, err
	}

	if receiverID != form.UserID{
		stmt = `
			INSERT INTO notifications (receiver_id, sender_id, author_name, action_type, snippet_id, comment_id, timestamp)
			VALUES (?, ?, ?, 'comment', ?, ?, datetime('now', '+6 hours'));
			`
		_, err = m.DB.Exec(stmt, receiverID, form.UserID, form.CommentatorName, form.SnippetID, int(commentID))
		if err != nil {
			return 0, err
		}
	}

	return int(commentID), nil
}

func (m *CommentModel) GetComments(snippetID int) ([]*Comment, error) {

	stmt := `
	SELECT comments.id, comments.snippet_id, comments.content, 
	comments.created, users.name, comments.user_id, comments.reported from comments
	JOIN users on users.id = comments.user_id
	WHERE comments.snippet_id = ?;
	`

	rows, err := m.DB.Query(stmt, snippetID)
	if err != nil {
		return nil, err
	}

	comments := []*Comment{}

	for rows.Next(){
		c := &Comment{}

		err = rows.Scan(&c.ID, &c.SnippetID, &c.Content, &c.Created, &c.Author, &c.AuthorID, &c.Reported)

		if err != nil{
			return nil, err
		}

		stmt = `
		SELECT COUNT(*) FROM reactions
		WHERE comment_id = ? AND type = 'likecomment';
		`
		row := m.DB.QueryRow(stmt, c.ID)

		err = row.Scan(&c.Likes)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Likes = 0
			} else {
				return nil, err
			}
		}

		stmt = `
		SELECT COUNT(*) FROM reactions
		WHERE comment_id = ? AND type = 'dislikecomment';
		`

		row = m.DB.QueryRow(stmt, c.ID)

		err = row.Scan(&c.Dislikes)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Dislikes = 0
			} else {
				return nil, err
			}
		}

		comments = append(comments, c)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (m *CommentModel) CommentReports() ([]*Comment, error){
	comments := []*Comment{}

	stmt := `
	SELECT comments.id, comments.snippet_id, comments.content, users.name, snippets.title FROM comments
	JOIN users ON comments.user_id = users.id
	JOIN snippets ON comments.snippet_id = snippets.id
	WHERE comments.reported = 1;
	`
	rows, err := m.DB.Query(stmt)
	if err != nil{
		return nil, err
	}

	for rows.Next() {
		c := &Comment{}

		err = rows.Scan(&c.ID, &c.SnippetID, &c.Content, &c.Author, &c.SnippetTitle)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return comments, nil

}

func (m *CommentModel) ReactedComments(userID int) ([]*Comment, error){
	comments := []*Comment{}
	stmt := `
	SELECT comments.id, comments.snippet_id, comments.content, comments.created, snippets.title, reactions.type FROM comments
	JOIN snippets on comments.snippet_id = snippets.id
	JOIN reactions on comments.id = reactions.comment_id
	WHERE reactions.user_id = ? AND (reactions.type = 'likecomment' OR reactions.type = 'dislikecomment')
	ORDER BY comments.created DESC;
	`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil{
		return nil, err
	}

	for rows.Next(){
		c := &Comment{}

		err = rows.Scan(&c.ID, &c.SnippetID, &c.Content, &c.Created, &c.SnippetTitle, &c.Status)
		if err != nil{
			return nil, err
		}

		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return comments, nil
	
}

func (m *CommentModel) MyComments(userID int) ([]*Comment, error){
	comments := []*Comment{}
	stmt := `
	SELECT comments.id, comments.snippet_id, comments.content, comments.created, snippets.title FROM comments
	JOIN snippets on comments.snippet_id = snippets.id
	WHERE comments.user_id = ?
	ORDER BY comments.created DESC;
	`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil{
		return nil, err
	}

	for rows.Next(){
		c := &Comment{}

		err = rows.Scan(&c.ID, &c.SnippetID, &c.Content, &c.Created, &c.SnippetTitle)
		if err != nil{
			return nil, err
		}

		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return comments, nil

}

func (m *CommentModel) GetComment(commentID int) (*Comment, error) {
	c := &Comment{}	
	
	stmt := `
	SELECT snippet_id, user_id, content FROM comments
	WHERE id = ?;
	`
	
	err := m.DB.QueryRow(stmt, commentID).Scan(&c.SnippetID, &c.AuthorID, &c.Content)
	if err != nil{
		return nil, err
	}

	c.ID = commentID
	return c, nil
}

func (m *CommentModel) EditComment(form *form.CommentEdit) error{
	stmt :=`
	UPDATE comments SET content = ?
	WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, form.Comment, form.CommentID)
	
	return err
}

func (m *CommentModel) DeleteComment(commentID int) error{

	stmt := `
	DELETE FROM comments WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, commentID)

	return err
}

func (m *CommentModel) ReportComment(form *form.Report) error{
	stmt := `
	UPDATE comments SET reported = 1
	WHERE id = ? AND snippet_id = ?;
	`

	_, err := m.DB.Exec(stmt, form.CommentID, form.SnippetID)
	return err
}

func (m *CommentModel) UnreportComment(commentID int) error{
	stmt := `
	UPDATE comments SET reported = 0
	WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, commentID)
	return err
}