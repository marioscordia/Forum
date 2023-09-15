package handler

import (
	"errors"
	"net/http"
	"newforum/internal/form"
	"newforum/internal/oops"
	"newforum/internal/temp"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
		case http.MethodGet:
			if tmpData.IsAuthenticated{
				h.ErrorHandler(w, http.StatusForbidden, tmpData)
				return
			}
			tmpData.Form = form.UserSignup{}
			h.render(w, http.StatusOK, "signup.html", tmpData)
		case http.MethodPost:
			if tmpData.IsAuthenticated{
				h.ErrorHandler(w, http.StatusForbidden, tmpData)
				return
			}
			err := r.ParseForm()
			if err != nil {
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
				return
			}

			form := form.UserSignup{
				Name: r.PostForm.Get("name"),
				Email: r.PostForm.Get("email"),
				Password: r.PostForm.Get("password"),
				Confirm: r.PostForm.Get("confirm"),
			}
			
			err = h.service.CreateUser(&form)
			if err != nil{
				if errors.Is(err, oops.ErrFormInvalid){
					tmpData.Form = form
					h.render(w, http.StatusUnprocessableEntity, "signup.html", tmpData)
				}else{
					h.Error(err)
					h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
				}
				return
			}

			http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		default:
			h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
  }
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
		case http.MethodGet:
			if tmpData.IsAuthenticated{
				h.ErrorHandler(w, http.StatusForbidden, tmpData)
				return
			}
			tmpData.Form = form.UserSignin{}
			h.render(w, http.StatusOK, "signin.html", tmpData)
		case http.MethodPost:
			if tmpData.IsAuthenticated{
				h.ErrorHandler(w, http.StatusForbidden, tmpData)
				return
			}
			
			err := r.ParseForm()
			if err != nil {
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
				return
			}

			form := form.UserSignin{
				Email: r.PostForm.Get("email"),
				Password: r.PostForm.Get("password"),
			}

			userID, err := h.service.Authenticate(&form)
    		if err != nil {
					if errors.Is(err, oops.ErrInvalidCredentials) {
						tmpData.Form = form
						h.render(w, http.StatusUnprocessableEntity, "signin.html", tmpData)
					} else {
						h.Error(err)
						h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
					} 
					return
    		}

			err = h.service.CreateSession(w, userID)
			if err != nil {
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		default:
			h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:

		cookie, err := r.Cookie("session")
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		err = h.service.DeleteSession(w, cookie)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) BeModerator(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:

		err := h.service.MakeRequest(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}