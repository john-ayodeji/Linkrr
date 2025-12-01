package utils

import (
	"bytes"
	"html/template"
	"log"
)

func RenderTemplate(filePath string, data interface{}) string {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Error executing template: %v", err)
	}

	return buf.String()
}
