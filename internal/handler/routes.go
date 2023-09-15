package handler

import "net/http"


func (h *Handler) Routes() http.Handler {
	
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	imgServer := http.FileServer(http.Dir("./internal/store/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/store/", http.StripPrefix("/store", imgServer))

	//General handlers
	mux.HandleFunc("/", h.Home)

	mux.HandleFunc("/filter", h.Category)

	mux.HandleFunc("/snippet/create", h.requireAuthentication(h.Create))

	mux.HandleFunc("/snippet/view", h.View)
	
	mux.HandleFunc("/snippet/created", h.requireAuthentication(h.Created))

	mux.HandleFunc("/snippet/commented", h.requireAuthentication(h.Commented))

	mux.HandleFunc("/reacted/snippet", h.requireAuthentication(h.ReactedSnippets))
	
	mux.HandleFunc("/reacted/comment", h.requireAuthentication(h.ReactedComments))

	mux.HandleFunc("/reaction/snippet", h.requireAuthentication(h.SnippetReaction))

	mux.HandleFunc("/reaction/comment", h.requireAuthentication(h.CommentReaction))

	mux.HandleFunc("/comment/create", h.requireAuthentication(h.CreateComment))

	mux.HandleFunc("/notifications", h.requireAuthentication(h.Notification))

	mux.HandleFunc("/snippet/edit", h.requireAuthentication(h.SnippetEdit))

	mux.HandleFunc("/comment/edit", h.requireAuthentication(h.CommentEdit))

	mux.HandleFunc("/snippet/delete", h.requireAuthentication(h.DeleteSnippet))

	mux.HandleFunc("/comment/delete", h.requireAuthentication(h.DeleteComment))

	mux.HandleFunc("/user/request", h.requireAuthentication(h.BeModerator))

	mux.HandleFunc("/user/signup", h.SignUp)

	mux.HandleFunc("/user/signin", h.SignIn)
	
	mux.HandleFunc("/user/logout", h.requireAuthentication(h.Logout))

	//Admin handlers
	mux.HandleFunc("/snippet/approve", h.OnlyAdmin(h.requireAuthentication(h.ApproveSnippet)))

	mux.HandleFunc("/comment/approve", h.OnlyAdmin(h.requireAuthentication(h.ApproveComment)))

	mux.HandleFunc("/user/requests", h.OnlyAdmin(h.requireAuthentication(h.Requests)))

	mux.HandleFunc("/user/list", h.OnlyAdmin(h.requireAuthentication(h.UserList)))

	mux.HandleFunc("/user/approval", h.OnlyAdmin(h.requireAuthentication(h.AdminApproval)))

	mux.HandleFunc("/user/role", h.OnlyAdmin(h.requireAuthentication(h.ChangeRole)))

	mux.HandleFunc("/report/snippet", h.OnlyAdmin(h.requireAuthentication(h.SnippetReports)))

	mux.HandleFunc("/report/comment", h.OnlyAdmin(h.requireAuthentication(h.CommentReports)))

	//Moderator handlers
	mux.HandleFunc("/pending", h.OnlyModerator(h.requireAuthentication(h.Pending)))

	mux.HandleFunc("/pending/view", h.OnlyModerator(h.requireAuthentication(h.PendingView)))

	mux.HandleFunc("/pending/approval", h.OnlyModerator(h.requireAuthentication(h.Approval)))
	
	mux.HandleFunc("/snippet/report", h.OnlyModerator(h.requireAuthentication(h.SnippetReport)))

	mux.HandleFunc("/comment/report", h.OnlyModerator(h.requireAuthentication(h.CommentReport)))

	return h.Middleware(h.logRequest(h.recoverPanic(mux)))
		
}