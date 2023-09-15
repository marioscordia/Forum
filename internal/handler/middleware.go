package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"newforum/internal/oops"
	"newforum/internal/temp"
	"time"
)

type contextKey string
const ctxKey contextKey = "data"

func (h *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := h.newTemplateData(r)
		cookie, err := r.Cookie("session")
		switchLabel:
		switch err {
			case nil:
				info, err := h.service.TakeInfo(cookie.Value)
				if err!=nil {
					if errors.Is(err, oops.ErrInvalidCredentials){
						data.IsAuthenticated = false
						data.UserInfo = temp.UserInfo{}
						cookie.Value = ""
						cookie.Path = "/"
						cookie.Expires = time.Now()
						cookie.MaxAge = -1 
						http.SetCookie(w, cookie)
						break switchLabel
					} else {
						h.Error(err)
						h.ErrorHandler(w, http.StatusInternalServerError, data)
						return
					}
				}
				
				data.UserInfo = temp.UserInfo{
					ID: info.ID,
					Name: info.Name,
					Role: info.Role,
					Requested: info.Requested,
				}
				data.IsAuthenticated = true
				nots, err := h.service.NotificationNum(data.ID)
				if err != nil {
					h.Error(err)
					h.ErrorHandler(w, http.StatusInternalServerError, data)
					return
				}
				data.NotNum = nots
			case http.ErrNoCookie:
				data.IsAuthenticated = false
				data.UserInfo = temp.UserInfo{}
			default:
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, data)
				return
		}
		ctx := context.WithValue(r.Context(), ctxKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) requireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(ctxKey).(*temp.TemplateData)
			// Otherwise set the "Cache-Control: no-store" header so that pages
			// require authentication are not stored in the users browser cache (or
			// other intermediary cache).
			if !data.IsAuthenticated{
				http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
				return
			}
			// w.Header().Add("Cache-Control", "no-store")
			// And call the next handler in the chain.
			next.ServeHTTP(w, r)
		})
}

func (h *Handler) OnlyAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(ctxKey).(*temp.TemplateData)
		if data.Role != 1{
			h.ErrorHandler(w, http.StatusForbidden, data)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func (h *Handler) OnlyModerator(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(ctxKey).(*temp.TemplateData)
		if data.Role != 2{
			h.ErrorHandler(w, http.StatusForbidden, data)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Create a deferred function (which will always be run in the event
	// of a panic as Go unwinds the stack).
		defer func() {
	// Use the builtin recover function to check if there has been a
	// panic or not. If there has...
			if err := recover(); err != nil {
	// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
	// Call the app.serverError helper method to return a 500
	// Internal Server response.
				fmt.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		h.infoLogger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}