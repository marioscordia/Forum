package handler

import (
	"fmt"
	"net/http"
	"newforum/internal/temp"
	"strconv"
)

func (h *Handler) DeleteSnippet(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)	
			return
		}

		id, err := strconv.Atoi(r.PostForm.Get("id"))
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)	
		}
		image := r.PostForm.Get("image")

		snippet, err := h.service.GetSnippet(id)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		if snippet.AuthorID != tmpData.ID{
			h.ErrorHandler(w, http.StatusUnauthorized, tmpData)
			return
		}

		err = h.service.DeleteSnippet(id, image)
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

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)	
			return
		}

		id := r.PostForm.Get("snippetid")
		
		cid, err := strconv.Atoi(r.PostForm.Get("commentid"))
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)	
		}

		comment, err := h.service.GetComment(cid)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}
		if comment.AuthorID != tmpData.ID{
			h.ErrorHandler(w, http.StatusUnauthorized, tmpData)
			return
		}

		err = h.service.DeleteComment(cid)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%s", id), http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}