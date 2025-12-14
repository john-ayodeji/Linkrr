package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

func RenderTemplate(filePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return "", fmt.Errorf("error parsing template")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Error executing template: %v", err)
		return "", fmt.Errorf("error executing template")
	}

	return buf.String(), nil
}
