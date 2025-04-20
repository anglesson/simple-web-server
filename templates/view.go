package templates

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	ErrorMessage string
}

func View(w http.ResponseWriter, templateName string, data any) {
	tmpl := template.Must(template.ParseFiles(
		"templates/layouts/base.html",
		"templates/pages/"+templateName+".html",
	))
	log.Printf("Data: %v", data)
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error executing template:", err)
		// Set the status code to 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		tmplError := template.Must(template.ParseFiles("templates/pages/500-error.html"))
		err = tmplError.Execute(w, PageData{ErrorMessage: "Internal Server Error"})
		if err != nil {
			log.Println("Error executing error template:", err)
			// If we can't even render the error page, we just write a plain text response
			w.Write([]byte("Internal Server Error"))
		}
		return
	}
}
