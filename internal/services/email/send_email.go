package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/john-ayodeji/Linkrr/internal/utils"
)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	fmt.Println(string(body))
}
