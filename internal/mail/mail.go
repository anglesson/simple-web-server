package mail

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
)

// Create body message with template
func NewEmail(templateName string, data any) string {
	templatePath := filepath.Join("internal", "mail", "templates", templateName+".html")
	baseTemplate := filepath.Join("internal", "mail", "templates", "template.html")

	tmpl, err := template.ParseFiles(baseTemplate, templatePath)
	if err != nil {
		log.Fatalf("failed to parse templates: %s", err)
	}

	var bodyString bytes.Buffer

	if err := tmpl.Execute(&bodyString, data); err != nil {
		log.Fatalf("failed to execute template: %s", err)
	}

	return bodyString.String()
}
