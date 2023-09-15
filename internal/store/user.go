package store

import (
	"database/sql"
	"errors"
	"newforum/internal/form"
	"newforum/internal/oops"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int
	Name string
	Email string
	Role int
	Requested int
	HashPassword []byte
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

func (m *UserModel) Insert(form *form.UserSignup) error {
	

	stmt := `INSERT OR IGNORE INTO users (name, email, hashed_password, created)
	          VALUES (?, ?, ?, datetime('now','+6 hours'));`

	res, err := m.DB.Exec(stmt, form.Name, form.Email, form.Password)
	if err != nil{
		return err
	}

	rowsAffected, err := res.RowsAffected();
	if err != nil{
		return err
	}

	if  int(rowsAffected) == 0{
		return oops.ErrDuplicateEmail
	}
	
	// var sqliteErr *sqlite3.Error
	// if err != nil {
	// 	if errors.As(err, &sqliteErr){
	// 		if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique{
	// 			if strings.Contains(sqliteErr.ExtendedCode.Error(), "users.email"){
	// 				return ErrDuplicateEmail
	// 			}
	// 		}
	// 	}
	// 	return err
	// }
	
	return nil
}

func (m *UserModel) Authenticate(form *form.UserSignin) (int, error) {

	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?;"
	err := m.DB.QueryRow(stmt, form.Email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, oops.ErrInvalidCredentials
		} else {
			return 0, err
		}
	} 

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(form.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, oops.ErrInvalidCredentials
		} else {
			return 0, err
		}
	} 
	
	return id, nil
}

func (m *UserModel) PutToken(userID int, token string) error{

	stmt := `
		UPDATE users
		SET token = ?
		WHERE id = ?;
	`
	_, err := m.DB.Exec(stmt, token, userID)
	
	if err != nil {
		return err
	}
	return nil

}

func (m *UserModel) TakeInfo(token string) (User, error) {
	var info User

	stmt := `
	SELECT id, name, role, requested FROM users
	WHERE token = ?;
	`
	err := m.DB.QueryRow(stmt, token).Scan(&info.ID, &info.Name, &info.Role, &info.Requested)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return User{}, oops.ErrInvalidCredentials
		}
		return User{}, err
	}

	return info, nil
}

func (m *UserModel) DeleteToken(token string) error {
	stmt :=
	`
	UPDATE users SET token = null WHERE token = ?;
	`
	_, err := m.DB.Exec(stmt, token)

	if err != nil {
		return err
	}

	return nil

}

func (m *UserModel) UserExists(email string) (bool, error) {
	stmt :=`
	SELECT 1 from users WHERE email = ?
	`
	var exists int
	err := m.DB.QueryRow(stmt, email).Scan(&exists)
	if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }

	return true, nil
}

func (m *UserModel) MakeRequest(userID int) error{
	stmt := `
	UPDATE users SET requested = 1
	WHERE id = ?
	`

	_, err := m.DB.Exec(stmt, userID)
	return err
}

func (m *UserModel) GetUsers() ([]*User, error) {
	users := []*User{}

	stmt := `
	SELECT id, name, role FROM users
	WHERE role!=1;
	`

	rows, err := m.DB.Query(stmt)
	if err != nil{
		return nil, err
	}
	
	for rows.Next(){
		u := &User{}

		err = rows.Scan(&u.ID, &u.Name, &u.Role)
		if err != nil{
			return nil, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	return users, nil
}

func (m *UserModel) GetRequests() ([]*User, error) {
	users := []*User{}

	stmt := `
	SELECT id, name from users WHERE requested = 1;
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := &User{}

		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}
	return users, nil
}

func (m *UserModel) Reject(userID int) error{
	stmt := `
	UPDATE users SET requested = 0 
	WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, userID)
	return err
}

func (m *UserModel) Upgrade(userID int) error {
	stmt := `
	UPDATE users SET role = 2, requested = 0
	WHERE id = ?;

	`

	_, err := m.DB.Exec(stmt, userID)
	return err
}

func (m *UserModel) Downgrade(userID int) error {
	stmt := `
	UPDATE users SET role = 3
	WHERE id = ?;
	`

	_, err := m.DB.Exec(stmt, userID)
	return err
}

