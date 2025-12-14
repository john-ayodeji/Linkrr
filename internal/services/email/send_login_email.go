package email

import (
	"fmt"

	"github.com/john-ayodeji/Linkrr/internal/utils"
)

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
	html, err := utils.RenderTemplate(Path+"login_email.html", data)
	if err != nil {
		return
	}
	sendEmail("Login", subject, text, html, name, email)
}
