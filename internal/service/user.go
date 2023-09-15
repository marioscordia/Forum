package service

import (
	"errors"
	"net/http"
	"newforum/internal/form"
	"newforum/internal/oops"
	"newforum/internal/store"
	"newforum/internal/validator"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	user store.UserStore
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{
		user: store.UserStore,
	}
}

func (s *UserService) CreateUser(form *form.UserSignup) error{
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.Confirm), "confirm", "This field cannot be blank")
	form.CheckField(validator.ConfirmPassword(form.Password, form.Confirm), "confirm", "Passwords do not match")


	if !form.Valid(){
		return oops.ErrFormInvalid
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
	if err != nil {
		return err
	}

	form.Password = string(hashedPassword)

	err = s.user.Insert(form)
	if err != nil {
		if errors.Is(err, oops.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			return oops.ErrDuplicateEmail
		} else {
			return err
		}
	}
	return nil
}


func (s *UserService) Authenticate(form *form.UserSignin) (int, error){
	
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid(){
		return 0, oops.ErrFormInvalid
	}

	id, err := s.user.Authenticate(form)
	if err != nil {
		if errors.Is(err, oops.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrectfd")
			return 0, oops.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	
	return id, nil
}


func (s *UserService) CreateSession(w http.ResponseWriter,userID int) error{
	token, err := uuid.NewV4()
	if err != nil {
		return err
	}

	err = s.user.PutToken(userID, token.String())
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: token.String(),
		Path: "/",
		Expires: time.Now().Add(15 * time.Minute),
	})

	return nil
}


func (s *UserService) DeleteSession(w http.ResponseWriter, cookie *http.Cookie) error{
	err := s.user.DeleteToken(cookie.Value)
	if err != nil {
		return err
	}

	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Now()
	cookie.MaxAge = -1 
	http.SetCookie(w, cookie)
	
	return nil
}

func (s *UserService) TakeInfo(token string) (store.User, error) {
	return s.user.TakeInfo(token)
}

func (s *UserService) MakeRequest(userID int) error{
	return s.user.MakeRequest(userID)
}
