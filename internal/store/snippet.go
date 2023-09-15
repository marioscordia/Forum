package store

import (
	"database/sql"
	"errors"
	"newforum/internal/form"
	"newforum/internal/oops"

	"time"
)

type Snippet struct {
	ID int
	Title string
	Content string
	Image string
	Created time.Time
	Author string
	AuthorID int
	Likes int
	Dislikes int
	Status string
	Approved int
	Reported int
	Category []string
}

type SnippetModel struct {
	DB *sql.DB
}

func NewSnippetModel(db *sql.DB) *SnippetModel{
	return &SnippetModel{DB: db}
}

func (m *SnippetModel) LatestSnippets() ([]*Snippet, error) {

	stmt := `
	SELECT snippets.id, snippets.title, snippets.created, users.name FROM snippets
	JOIN users ON snippets.user_id = users.id
	WHERE snippets.approved = 1
	ORDER BY snippets.id DESC;	
	`
	
	rows, err := m.DB.Query(stmt)
	
	if err != nil {
		return nil, err
	}

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Author)
		
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) SnippetReports() ([]*Snippet, error){
	snippets := []*Snippet{}
	stmt := `
	SELECT snippets.id, snippets.title, users.name from snippets
	JOIN users on snippets.user_id = users.id
	WHERE snippets.reported = 1; 
	`
	rows, err := m.DB.Query(stmt)
	if err != nil{
		return nil, err
	}

	for rows.Next(){
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Author)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil{
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) PendingSnippets() ([]*Snippet, error) {
	stmt := `
	SELECT snippets.id, snippets.title, snippets.created, users.name FROM snippets
	JOIN users ON snippets.user_id = users.id
	WHERE snippets.approved = 0
	ORDER BY snippets.id DESC;	
	`
	
	rows, err := m.DB.Query(stmt)
	
	if err != nil {
		return nil, err
	}

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Author)
		
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) CreateSnippet(form *form.SnippetCreate) (int, error) {

	stmt := `
	INSERT INTO snippets (user_id, title, content, image, created)
	VALUES (?, ?, ?, ?, datetime('now','+6 hours'));
	`

	result, err := m.DB.Exec(stmt, form.UserID, form.Title, form.Content, form.FileName)
	if err != nil {
		return 0, err
	} 

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	stmt = `
	INSERT INTO categories (snippet_id, category)
	VALUES (?, ?);
	`

	for _, c := range form.Category{
		_, err = m.DB.Exec(stmt, int(id), c)
		if err != nil {
			return 0, err
		}
	}

	return int(id), nil
} 

func (m *SnippetModel) GetSnippet(snippetID int) (*Snippet, error) {
	stmt := `SELECT snippets.id, snippets.title, snippets.content, 
					snippets.image, snippets.created, users.name, 
					snippets.user_id, snippets.approved, snippets.reported FROM snippets
					JOIN users ON snippets.user_id = users.id 			
					WHERE snippets.id = ?;`

	row := m.DB.QueryRow(stmt, snippetID)
	
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Image, &s.Created, &s.Author, &s.AuthorID, &s.Approved, &s.Reported)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, oops.ErrNoRecord
		} else {
			return nil, err
		}
	}

	categories := []string{}
	
	stmt = `SELECT category FROM categories
			WHERE snippet_id = ?`
	
	rows, err := m.DB.Query(stmt, snippetID)
	
	if err != nil {
		return nil, err
	}

	for rows.Next(){
		var c string

		err = rows.Scan(&c)
		if err != nil{
			return nil, err
		}

		categories = append(categories, c)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}

	s.Category = categories	

	stmt = `
	SELECT COUNT(*) FROM reactions
	WHERE snippet_id = ? AND type = 'likepost';
	`
	row = m.DB.QueryRow(stmt, snippetID)

	err = row.Scan(&s.Likes)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Likes = 0
		} else {
			return nil, err
		}
	}

	stmt = `
	SELECT COUNT(*) FROM reactions
	WHERE snippet_id = ? AND type = 'dislikepost';
	`

	row = m.DB.QueryRow(stmt, snippetID)

	err = row.Scan(&s.Dislikes)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Dislikes = 0
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) FilterSnippets(filter form.Filter) ([]*Snippet, error){
	snippets := []*Snippet{}

	stmt := `
	SELECT snippets.id, snippets.title, snippets.created, users.name 
	FROM snippets
	JOIN categories ON snippets.id = categories.snippet_id
	JOIN users on snippets.user_id = users.id  
	WHERE categories.category = ?
	ORDER BY snippets.id DESC;
	`

	for _, v := range filter.Category{
	
		rows, err := m.DB.Query(stmt, v)
	
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			s := &Snippet{}
			
			err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Author)
			
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				} else {
					return nil, err
				}
			} 

			if !Contains(snippets, s){
				snippets = append(snippets, s)
			}
		}

		if err= rows.Err(); err != nil {
			return nil, err
		}
	}

	return snippets, nil
}

func (m *SnippetModel) CreatedSnippets(userID int) ([]*Snippet, error){
	snippets := []*Snippet{}
	stmt:=`
		SELECT snippets.id, snippets.title, snippets.created, users.name
		FROM snippets
		JOIN users ON snippets.user_id = users.id
		WHERE snippets.user_id = ?
		ORDER BY snippets.id DESC;
		`
	rows, err := m.DB.Query(stmt, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := &Snippet{}
		
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Author)
		
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, oops.ErrNoRecord
			} else {
				return nil, err
			}
		} 

		snippets = append(snippets, s)
		
	}

	if err= rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) ReactedSnippets(userID int) ([]*Snippet, error){
	snippets := []*Snippet{}	

	stmt := `
	SELECT snippets.id, snippets.title, snippets.created, users.name, reactions.type
	FROM snippets
	JOIN reactions ON snippets.id = reactions.snippet_id
	JOIN users on snippets.user_id = users.id
	WHERE reactions.user_id = ? AND (reactions.type = 'likepost' OR reactions.type = 'dislikepost')
	ORDER BY snippets.id DESC;
	`
	rows, err := m.DB.Query(stmt, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := &Snippet{}
		
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Author, &s.Status)
		
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, oops.ErrNoRecord
			} else {
				return nil, err
			}
		} 

		snippets = append(snippets, s)
		
	}

	if err= rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}	

func (m *SnippetModel) UpdateSnippet(form *form.SnippetEdit) error{
	stmt := `
	UPDATE snippets SET title = ?, content = ?
	WHERE id = ?;
	`
	_, err := m.DB.Exec(stmt, form.Title, form.Content, form.SnippetID)
	
	return err
}

func (m *SnippetModel) DeleteSnippet(snippetID int) error{

	stmt := `
	DELETE FROM snippets WHERE id = ?;
	`
	_, err := m.DB.Exec(stmt, snippetID)

	return err
}

func (m *SnippetModel) ApproveSnippet(snippetID int) error{
	stmt := `
	UPDATE snippets SET approved = 1
	WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, snippetID)
	return err
}

func (m *SnippetModel) UnreportSnippet(snippetID int) error{
	stmt := `
	UPDATE snippets SET reported = 0
	WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, snippetID)
	return err
}

func (m *SnippetModel) ReportSnippet(snippetID int) error{
	stmt := `
	UPDATE snippets SET reported = 1
	WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, snippetID)
	return err
}

func Contains(snippets []*Snippet, s *Snippet) bool {

	if len(snippets) == 0 {
		return false
	}

	for _, v := range snippets{
		if v.ID == s.ID{
			return true
		}
	}

	return false
}
