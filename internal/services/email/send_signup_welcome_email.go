package email

import (
	"fmt"

	"github.com/john-ayodeji/Linkrr/internal/utils"
)

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
	html := utils.RenderTemplate("internal/email_templates/signup_email.html", data)
	sendEmail("Signup", subject, text, html, name, email)
}
