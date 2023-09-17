package handler

import (
	"errors"
	"fmt"
	"net/http"
	"newforum/internal/form"
	"newforum/internal/oops"
	"newforum/internal/temp"
	"strconv"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:
		tmpData.Form = form.SnippetCreate{}
		h.render(w, http.StatusOK, "create.html", tmpData)
	case http.MethodPost:
		if err := r.ParseMultipartForm(1024 * 1024 * 20); err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		file, handler, err := r.FormFile("image_name")
		if err != nil {
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
				return		
		}
		defer file.Close()
		
		form := form.SnippetCreate{
			UserID: tmpData.ID,
			Title: r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
			Category: r.PostForm["category"],
			FileName: handler.Filename,
			FileSize: int(handler.Size),
			File: file,
		}

		_, err = h.service.CreateSnippet(&form)
		if err != nil {
			if errors.Is(err, oops.ErrFormInvalid){
				tmpData.Form = form
				h.render(w, http.StatusBadRequest, "create.html", tmpData)
			}else{
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			}
			return
		}
	
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}


func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			h.Error(err)
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}

		snippet, err := h.service.GetSnippet(id)
		if err != nil {
				if errors.Is(err, oops.ErrNoRecord) {
					h.ErrorHandler(w, http.StatusNotFound, tmpData)
				} else {
					h.Error(err)
					h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
				}
				return 
		}
		
		comments, err := h.service.GetComments(id)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Snippet = snippet
		tmpData.Comments = comments
		tmpData.Form = form.Comment{}

		switch tmpData.Role{
		case 1:
			h.render(w, http.StatusOK, "adminview.html", tmpData)
		case 2:
			h.render(w, http.StatusOK, "modview.html", tmpData)
		default:
			h.render(w, http.StatusOK, "view.html", tmpData)
		}		
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:
		snippetID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || snippetID < 1 {
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

		form := form.Comment{
			UserID: tmpData.ID,
			CommentatorName: tmpData.Name,
			SnippetID: snippetID,
			Comment: r.PostForm.Get("comment"),
		}

		commentID, err := h.service.CreateComment(&form)
		if err != nil{
			if errors.Is(err, oops.ErrFormInvalid){
				url := fmt.Sprintf("/snippet/view?id=%d", snippetID)
				http.Redirect(w, r, url, http.StatusSeeOther)
				// tmpData.Form = form
				// ctx := context.WithValue(r.Context(), ctxKey, tmpData)
				// http.Redirect(w, r.WithContext(ctx), fmt.Sprintf("/snippet/view?id=%d#form", snippetID), http.StatusSeeOther)
				
				// // h.render(w, http.StatusUnprocessableEntity, "view.html", tmpData)
			}else{
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)		
			}
			return
		}
		url := fmt.Sprintf("/snippet/view?id=%d#%d", snippetID, commentID)
		http.Redirect(w, r, url, http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}