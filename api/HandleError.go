package api

import (
	"html/template"
	"log"
	"net/http"
)

func responseAlreadyWritten(w http.ResponseWriter) bool {
	return w.Header().Get("Content-Type") != ""
}

// HandleError handles HTTP errors, attempting to display a specific error page.
// If the specific page (like "404.html") is unavailable, it falls back to "500.html".
func HandleError(w http.ResponseWriter, err error, statusCode int, templateName string) {
	
	// Attempt to parse the specified error template
	tpl, tplErr := template.ParseFiles("templates/" + templateName)
	if tplErr != nil {
		log.Println("Error parsing template", templateName, ":", tplErr)

		// If the specified template cannot be parsed, set status to 500 and use the fallback
		if !responseAlreadyWritten(w) {
			w.WriteHeader(http.StatusInternalServerError)
		}
		tpl, _ = template.ParseFiles("templates/500.html") // Attempt to load the 500 fallback template
	}

	
	tpl.Execute(w, nil)
}
