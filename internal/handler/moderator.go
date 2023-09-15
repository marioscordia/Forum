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

func (h *Handler) Pending(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:

		snippets, err := h.service.PendingSnippets()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Snippets = snippets

		h.render(w, http.StatusOK, "pending.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) PendingView(w http.ResponseWriter, r *http.Request) {
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

		tmpData.Snippet = snippet

		h.render(w, http.StatusOK, "pendingview.html", tmpData)	
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) Approval(w http.ResponseWriter, r*http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		id, err := strconv.Atoi(r.PostForm.Get("id"))
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.Approval{
			SnippetID: id,
			Image: r.PostForm.Get("image"),
			Decision: r.PostForm.Get("decision"),
		}

		err = h.service.Approval(&form)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, "/pending", http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) SnippetReport(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodPost:
	
		if err := r.ParseForm(); err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		snippetID, err := strconv.Atoi(r.PostForm.Get("id"))
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		err = h.service.ReportSnippet(snippetID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", snippetID), http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) CommentReport(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		snippetID, err := strconv.Atoi(r.PostForm.Get("snippetid"))
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		commentID, err := strconv.Atoi(r.PostForm.Get("commentid"))
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.Report{
			SnippetID: snippetID,
			CommentID: commentID,
		}

		err = h.service.ReportComment(&form)
		if err != nil { 
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d#%d", snippetID, commentID), http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}
