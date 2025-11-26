package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/john-ayodeji/Linkrr/utils"
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

func SendWelcomeEmail(name, email string) {
	subject := fmt.Sprintf("Welcome to Linkrr, %s!", name)
	text := fmt.Sprintf("Hi %s, welcome to Linkrr! We're excited to have you aboard.", name)
	data := struct {
		Name     string
		LoginURL string
	}{
		Name:     name,
		LoginURL: "localhost:3000/api/v1/auth/login",
	}
	html := utils.RenderTemplate("./email_templates/signup_email.html", data)
	sendEmail("Signup", subject, text, html, name, email)
}

func SendLoginWelcomeEmail(name, email string) {
	subject := fmt.Sprintf("Welcome back, %s!", name)
	text := fmt.Sprintf("Hi %s, great to see you again on Linkrr.", name)
	data := struct {
		Name     string
		LoginURL string
	}{
		Name:     name,
		LoginURL: "localhost:3000/api/v1/auth/login",
	}
	html := utils.RenderTemplate("./email_templates/login_email.html", data)
	sendEmail("Login", subject, text, html, name, email)
}

func SendPasswordResetEmail(name, email, url string) {
	subject := fmt.Sprintf("Password Reset Request")
	text := fmt.Sprintf("Hello %v\n We received a request to reset your password.\nClick the link below to set a new password\n%v\nThis link will expire in 15 minutes.\n\n\nIf you did not request this, you can safely ignore this email\nThanks, The Linkrr Team", name, url)
	data := struct {
		Name string
		URL  string
	}{
		Name: name,
		URL:  url,
	}
	html := utils.RenderTemplate("./email_templates/password_email.html", data)
	sendEmail("Password-Reset", subject, text, html, name, email)
}
