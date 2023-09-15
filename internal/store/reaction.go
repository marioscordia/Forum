package store

import (
	"database/sql"
	"newforum/internal/form"
)

type ReactionModel struct {
	DB *sql.DB
}

func NewReactionModel(db *sql.DB) *ReactionModel {
	return &ReactionModel{DB: db}
}

var mapReaction = map[string]string{
	"likepost": "dislikepost",
	"dislikepost": "likepost",
	"likecomment": "dislikecomment",
	"dislikecomment": "likecomment",
}


//  This is the part for post reactions

func (m *ReactionModel) PostReaction(form form.SnippetReaction) error {

	exists, err := m.ExistsReactionPost(form.SnippetID, form.UserID, form.Reaction)

	if err != nil{
		return err
	}

	if exists{
		err = m.DeleteReactionPost(form.SnippetID, form.UserID, form.Reaction)
		if err != nil{
			return err
		}
	}else{
		stmt := `
		INSERT OR IGNORE INTO reactions (snippet_id, user_id, type)
		VALUES (?, ?, ?);
		`
		
		_, err = m.DB.Exec(stmt, form.SnippetID, form.UserID, form.Reaction)
	
		if err != nil{
			return err
		}
	
		exists, err = m.ExistsReactionPost(form.SnippetID, form.UserID, mapReaction[form.Reaction])

		if err != nil{
			return err
		}

		if exists{
			err = m.DeleteReactionPost(form.SnippetID, form.UserID, mapReaction[form.Reaction])
			if err != nil{
				return err
			}
		}

		var receiverID int

		stmt = `
		SELECT user_id FROM snippets
		WHERE id = ?;
		`
		err = m.DB.QueryRow(stmt, form.SnippetID).Scan(&receiverID)
		if err != nil {
			return err
		}

		if receiverID != form.UserID{
			stmt = `
			INSERT INTO notifications (receiver_id, sender_id, author_name, action_type, snippet_id, timestamp)
			VALUES (?, ?, ?, ?, ?, datetime('now', '+6 hours'));
			`
			_, err = m.DB.Exec(stmt, receiverID, form.UserID, form.ReactorName, form.Reaction, form.SnippetID)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (m *ReactionModel) DeleteReactionPost(snippetID, userID int, reaction string) error{
	stmt :=`
	DELETE FROM reactions 
	WHERE snippet_id = ? AND user_id = ? AND type = ?;
	`
	_, err := m.DB.Exec(stmt, snippetID, userID, reaction)

	if err != nil{
		return err
	}
	return nil

}

func (m *ReactionModel) ExistsReactionPost(snippetID, userID int, reaction string) (bool, error){
	stmt := `
	SELECT EXISTS(SELECT 1 FROM reactions WHERE snippet_id = ? AND user_id = ? AND type = ?);
	`
	var exists int
	err := m.DB.QueryRow(stmt, snippetID, userID, reaction).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
// ____________________________________________________________________________________________________________________

// This is the part for comment reactions

func (m *ReactionModel) CommentReaction(form form.CommentReaction) error {


	exists, err := m.ExistsReactionComment(form.CommentID, form.UserID, form.Reaction)

	if err != nil{
		return err
	}

	if exists{
		err = m.DeleteReactionComment(form.CommentID, form.UserID, form.Reaction)
		if err != nil{
			return err
		}
	}else{

		stmt := `
		INSERT OR IGNORE INTO reactions (comment_id, user_id, type)
		VALUES (?, ?, ?);
		`
		_, err = m.DB.Exec(stmt, form.CommentID, form.UserID, form.Reaction)

		if err != nil{
			return err
		}

		exists, err = m.ExistsReactionComment(form.CommentID, form.UserID, mapReaction[form.Reaction])

		if err != nil{
			return err
		}

		if exists{
			err = m.DeleteReactionComment(form.CommentID, form.UserID, mapReaction[form.Reaction])
			if err != nil{
				return err
			}
		}
	}
	return nil
	
}
	

func (m *ReactionModel) DeleteReactionComment(commentID, userID int, reaction string) error{
	stmt :=`
	DELETE FROM reactions 
	WHERE comment_id = ? AND user_id = ? AND type = ?;
	`
	_, err := m.DB.Exec(stmt, commentID, userID, reaction)

	if err != nil{
		return err
	}
	return nil

}

func (m *ReactionModel) ExistsReactionComment(commentID, userID int, reaction string) (bool, error){
	stmt := `
	SELECT EXISTS(SELECT 1 FROM reactions WHERE comment_id = ? AND user_id = ? AND type = ?);
	`
	var exists int
	err := m.DB.QueryRow(stmt, commentID, userID, reaction).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
// ___________________________________________________________________________________________________________

