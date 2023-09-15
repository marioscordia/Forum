package handler

import (
	"net/http"
	"newforum/internal/form"
	"newforum/internal/temp"
	"newforum/internal/validator"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:

		if r.URL.Path != "/" {
			h.ErrorHandler(w, http.StatusNotFound, tmpData)
			return
		}

		snippets, err := h.service.GetSnippets()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.Snippets = snippets
		tmpData.Form = form.Filter{}

		h.render(w, http.StatusOK, "home.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}

}

func (h *Handler) Category(w http.ResponseWriter, r *http.Request) {
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method{
	case http.MethodGet:
		err := r.ParseForm()
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		form := form.Filter{
			Category: r.URL.Query()["category"],
		}

		form.CheckField(validator.CheckCategory(form.Category), "filter", "Choose filter")

		if !form.Valid() {
			tmpData.Form = form
			h.render(w, http.StatusBadRequest, "home.html", tmpData)
			return
		}

		snippets, err := h.service.FilterSnippets(form)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}
		
		tmpData.Snippets = snippets

		h.render(w, http.StatusAccepted, "filter.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}