package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"newforum/internal/temp"
	"runtime/debug"
	"time"
)

func (h *Handler) render(w http.ResponseWriter, status int, page string, data *temp.TemplateData) {

	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} 
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (h *Handler) newTemplateData(r *http.Request) *temp.TemplateData {
	return &temp.TemplateData{
		CurrentYear: time.Now().Year(),
	}
}

func (h *Handler) Error(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.errorLogger.Output(2, trace)
}

// func getGithubClientID() string {

//     githubClientID, exists := os.LookupEnv("GITHUB_ID")
//     if !exists {
//         log.Fatal("Github Client ID not defined in .env file")
//     }

//     return githubClientID
// }

// func getGithubClientSecret() string {

//     githubClientSecret, exists := os.LookupEnv("GITHUB_SECRET")
//     if !exists {
//         log.Fatal("Github Client ID not defined in .env file")
//     }

//     return githubClientSecret
// }

// func parseEnvFile(path string) (map[string]string, error) {
// 	envMap := make(map[string]string)

// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())
// 		if len(line) > 0 && !strings.HasPrefix(line, "#") {
// 			parts := strings.SplitN(line, "=", 2)
// 			if len(parts) == 2 {
// 				key := strings.TrimSpace(parts[0])
// 				value := strings.TrimSpace(parts[1])
// 				envMap[key] = value
// 			}
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return nil, err
// 	}

// 	return envMap, nil
// }