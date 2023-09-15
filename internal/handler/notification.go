package handler

import (
	"net/http"
	"newforum/internal/temp"
)

func (h *Handler) Notification(w http.ResponseWriter, r *http.Request){
	tmpData := r.Context().Value(ctxKey).(*temp.TemplateData)
	switch r.Method {
	case http.MethodGet:
		nots, err := h.service.Notifications(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}
		tmpData.Notifications = nots
		
		err = h.service.Update(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		tmpData.NotNum, err = h.service.NotificationNum(tmpData.ID)
		if err != nil {
			h.Error(err)
			h.ErrorHandler(w, http.StatusInternalServerError, tmpData)
			return
		}

		h.render(w, http.StatusOK, "notifications.html", tmpData)
	default:
		h.ErrorHandler(w, http.StatusMethodNotAllowed, tmpData)
	}

}