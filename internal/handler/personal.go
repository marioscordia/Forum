package handler

import (
	"net/http"
	"newforum/internal/temp"
)

func (h *Handler) Created(w http.ResponseWriter, r *http.Request){
	
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		snippets, err := h.service.GetCreated(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Snippets = snippets
		h.render(w, http.StatusAccepted, "created.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) Commented(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		comments, err := h.service.MyComments(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Comments = comments
		h.render(w, http.StatusAccepted, "commented.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) ReactedSnippets(w http.ResponseWriter, r *http.Request){
	
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		snippets, err := h.service.ReactedSnippets(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Snippets = snippets
		h.render(w, http.StatusAccepted, "reactedsnipps.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}

func (h *Handler) ReactedComments(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		comments, err := h.service.ReactedComments(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Comments = comments
		h.render(w, http.StatusAccepted, "reactedcomms.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}