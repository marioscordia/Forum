package handler

import (
	"fmt"
	"net/http"
	"newforum/internal/form"
	"newforum/internal/temp"
	"strconv"
)

func (h *Handler) Requests(w http.ResponseWriter, r*http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:
		users, err := h.service.GetRequests()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Users = users
		h.render(w, http.StatusOK, "requests.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) AdminApproval(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		id, err := strconv.Atoi(r.PostForm.Get("id"))
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.AdminApproval{
			UserID: id,
			Decision: r.PostForm.Get("decision"),
		}

		if err = h.service.AdminApproval(&form); err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, "/user/requests", http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) UserList(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:
		users, err := h.service.GetUsers()
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Users = users
		h.render(w, http.StatusOK, "users.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) ChangeRole(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		id, err := strconv.Atoi(r.PostForm.Get("id"))
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.AdminApproval{
			UserID: id,
			Decision: r.PostForm.Get("decision"),
		}

		if err = h.service.AdminApproval(&form); err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, "/user/list", http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) ApproveSnippet(w http.ResponseWriter, r*http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:
		if err := r.ParseForm(); err != nil{
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

		if err = h.service.UnreportSnippet(id); err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) ApproveComment(w http.ResponseWriter, r*http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		snippetID := r.PostForm.Get("snippetid")
		commentID, err := strconv.Atoi(r.PostForm.Get("commentid"))
		if err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		if err = h.service.UnreportComment(commentID); err != nil{
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%s#%d", snippetID, commentID), http.StatusSeeOther)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) SnippetReports(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		snippets, err := h.service.SnippetReports()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Snippets = snippets
		h.render(w, http.StatusOK, "reports1.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) CommentReports(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		comments, err := h.service.CommentReports()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}
		
		tmpData.Comments = comments
		h.render(w, http.StatusOK, "reports2.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}