package handler

import (
	"net/http"
	"newforum/internal/temp"
)

// var (
// 	ErrFooNotFound = errors.New("page not found")
// 	ErrFooMethodNotAllowed = errors.New("method not allowed")
// 	ErrFooUnauthorized = errors.New("unauthorized request")
// 	ErrFooNotAuthorized = errors.New("unauthorized request")
// )

func (h *Handler) ErrorHandler(w http.ResponseWriter, status int, data *temp.TemplateData) {

	errorForm := temp.ErrorInfo{
		Code: status,
		Text:  http.StatusText(status),
	}

	data.ErrorInfo = errorForm

	h.render(w, errorForm.Code, "error.html", data)
}