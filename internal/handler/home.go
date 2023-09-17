package handler

import (
	"errors"
	"net/http"
	"newforum/internal/form"
	"newforum/internal/oops"
	"newforum/internal/temp"
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

		snippets, err := h.service.FilterSnippets(&form)
		if err != nil {
			if errors.Is(err, oops.ErrFormInvalid){
				tmpData.Form = form
				h.render(w, http.StatusUnprocessableEntity, "home.html", tmpData)
			}else{
				h.Error(err)
				h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			}
			return
		}
		
		tmpData.Snippets = snippets

		h.render(w, http.StatusAccepted, "filter.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}
}