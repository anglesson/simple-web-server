package template

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	ErrorMessage string
}

func View(w http.ResponseWriter, templateName string, data any, layout string) {
	files := []string{
		fmt.Sprintf("web/templates/layouts/%s.html", layout),
		fmt.Sprintf("web/templates/pages/%s.html", templateName),
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Erro ao renderizar template", http.StatusInternalServerError)
		return
	}
}
