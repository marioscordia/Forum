package handler

import (
	"errors"
	"fmt"
	"net/http"
	"newforum/internal/form"
	"newforum/internal/oops"
	"newforum/internal/store"
	"newforum/internal/temp"
	"strconv"
)

func (h *Handler) SnippetEdit(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}

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

		form := form.SnippetEdit{
			Title: snippet.Title,
			Content: snippet.Content,
		}

		tmpData.Form = form
		tmpData.Snippet = snippet

		h.render(w, http.StatusOK, "editsnippet.html", tmpData)
	case http.MethodPost:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}

		if err = r.ParseForm(); err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.SnippetEdit{
			SnippetID: id,
			Title: r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
		}

		if err = h.service.EditSnippet(&form); err != nil{
			if errors.Is(err, oops.ErrFormInvalid){
				tmpData.Form = form
				tmpData.Snippet = &store.Snippet{ID: id}
				h.render(w, http.StatusUnprocessableEntity, "editsnippet.html", tmpData)
			}else{
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			}
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}


func (h *Handler) CommentEdit(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}

		cid, err := strconv.Atoi(r.URL.Query().Get("cid"))
		if err != nil || id < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
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
		
		tmpData.Comment = &store.Comment{SnippetID: id, ID: cid}
		tmpData.Form = form.CommentEdit{
			Comment: comment.Content,
		}
		h.render(w, http.StatusOK, "editcomment.html", tmpData)
	case http.MethodPost:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}

		cid, err := strconv.Atoi(r.URL.Query().Get("cid"))
		if err != nil || id < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}
		
		if err = r.ParseForm(); err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.CommentEdit{
			CommentID: cid,
			Comment: r.PostForm.Get("comment"),
		}

		if err = h.service.EditComment(&form); err != nil{
			if errors.Is(err, oops.ErrFormInvalid){
				tmpData.Comment = &store.Comment{SnippetID: id, ID: cid}
				tmpData.Form = form
				h.render(w, http.StatusUnprocessableEntity, "editcomment.html", tmpData)
			}else{
				h.Error(err)	
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			}
			return
		}

		url := fmt.Sprintf("/snippet/view?id=%d#%d", id, cid)
		http.Redirect(w, r, url, http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

