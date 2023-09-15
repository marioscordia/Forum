package store

import (
	"database/sql"
	"time"
)

type Notification struct {
	Author string
	Action string
	SnippetID int
	CommentID int
	Timestamp time.Time
}

type Notifications struct {
	DB *sql.DB
}

func NewNotificationsModel(db *sql.DB) *Notifications{
	return &Notifications{DB: db}
}

func (m *Notifications) Notifications(userID int)([]*Notification, error) {
	var notifications []*Notification

	stmt := `
	SELECT author_name, action_type, snippet_id, comment_id, timestamp FROM notifications
	WHERE receiver_id = ? AND sender_id != ?
	ORDER BY timestamp DESC;
	`
	rows, err := m.DB.Query(stmt, userID, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		n := &Notification{}

		err = rows.Scan(&n.Author, &n.Action, &n.SnippetID, &n.CommentID, &n.Timestamp)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, n)

	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return notifications, nil

}

func (m *Notifications) Update(userID int) error {
	stmt := `
	UPDATE notifications SET is_read = 1
	WHERE receiver_id = ?;
	`

	_, err := m.DB.Exec(stmt, userID)
	
	return err
}

func (m *Notifications) NotificationNum(userID int) (int, error){
	var num int

	stmt := `
	SELECT COUNT(*) FROM notifications
	WHERE receiver_id = ? AND sender_id != ? AND is_read = 0;
	`
	err := m.DB.QueryRow(stmt, userID, userID).Scan(&num)
	if err != nil{
		return 0, err
	}

	return num, nil
}