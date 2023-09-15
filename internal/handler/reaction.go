package handler

import (
	"fmt"
	"net/http"
	"newforum/internal/form"
	"newforum/internal/temp"
	"strconv"
)

func (h *Handler) SnippetReaction(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:
		snippetID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}
		err = r.ParseForm()
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.SnippetReaction{
			UserID: tmpData.ID,
			ReactorName: tmpData.Name,
			SnippetID: snippetID,
			Reaction: r.PostForm.Get("reaction"),
		}

		err = h.service.SnippetReaction(form)
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}
		
		url := fmt.Sprintf("/snippet/view?id=%d", snippetID)
		http.Redirect(w, r, url, http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) CommentReaction(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)

	switch r.Method{
	case http.MethodPost:

		snippetID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || snippetID < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}
		commentID, err := strconv.Atoi(r.URL.Query().Get("cid"))
		if err != nil || snippetID < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		err = r.ParseForm()		
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.CommentReaction{
			UserID: tmpData.ID,
			CommentID: commentID,
			Reaction: r.PostForm.Get("reaction"),
		}

		err = h.service.CommentReaction(form)
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
		}
		
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d#%d", snippetID, commentID), http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}