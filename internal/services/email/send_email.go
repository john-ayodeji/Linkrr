package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/john-ayodeji/Linkrr/internal/utils"
)

func GetEmailTemplateFolder() string {
	var platform = os.Getenv("PLATFORM")
	if platform == "local" || platform == "dev" {
		log.Println("internal/email_templates/")
		return "internal/email_templates/"
	}
	if platform == "docker" || platform == "prod" {
		return "/app/internal/email_templates/"
	}

	log.Println("internal/email_templates/")
	return "internal/email_templates/"
}

var Path = GetEmailTemplateFolder()

type mailtrapEmail struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type mailtrapPayload struct {
	From     mailtrapEmail   `json:"from"`
	To       []mailtrapEmail `json:"to"`
	Subject  string          `json:"subject"`
	Text     string
	HTML     string
	Category string `json:"category"`
}

func sendEmail(category, subject, text, html, name, email string) {
	token := os.Getenv("MAILTRAP_TOKEN")
	if token == "" {
		utils.LogError("MAILTRAP_TOKEN not set")
		return
	}

	payload := mailtrapPayload{
		From:     mailtrapEmail{Email: "linkrr@johnayodeji.dev", Name: "Linkrr"},
		To:       []mailtrapEmail{{Email: email, Name: name}},
		Subject:  subject,
		Text:     text,
		HTML:     html,
		Category: category,
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	req, err := http.NewRequest("POST", "https://send.api.mailtrap.io/api/send", bytes.NewReader(buf))
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	defer res.Body.Close()

	_, err1 := io.ReadAll(res.Body)
	if err1 != nil {
		utils.LogError(err1.Error())
		return
	}
}
