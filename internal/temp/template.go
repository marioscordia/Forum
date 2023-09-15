package temp

import (
	"html/template"
	"newforum/internal/store"
	"path/filepath"
	"time"
)


type TemplateData struct {
	UserInfo
	CurrentYear int
	Snippet *store.Snippet
	Snippets []*store.Snippet
	Comment *store.Comment
	Comments []*store.Comment
	Notifications []*store.Notification
	Users []*store.User
	NotNum int
	Form any
	IsAuthenticated bool
	ErrorInfo
}

type UserInfo struct{
	ID int
	Name string
	Role int
	Requested int
}

type ErrorInfo struct {
	Code int
	Text string
}

func HumanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func ShortenStr(str string) string {
	if len(str) > 10 {
		return str[:10]+"..."
	}
	return str
}

var functions = template.FuncMap{
	"humanDate": HumanDate,
	"shortenStr": ShortenStr,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	} 
	for _, page := range pages {
		
		name := filepath.Base(page)
		
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		} 

		ts, err = ts.ParseGlob("./ui/html/partial/*.html")
		if err != nil {
			return nil, err
		} 

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	} 

	return cache, nil
}