package handler

import (
	"html/template"
	"log"
	"newforum/internal/service"
)

type Handler struct {
	infoLogger *log.Logger
	errorLogger *log.Logger
	templateCache map[string]*template.Template
	service *service.Service
}

func NewHandler(info, err *log.Logger, cache map[string]*template.Template, service *service.Service) *Handler{
	return &Handler{
		infoLogger: info,
		errorLogger: err,
		templateCache: cache,
		service: service,
	}
}

