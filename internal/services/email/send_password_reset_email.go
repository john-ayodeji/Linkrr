package email

import (
	"fmt"

	"github.com/john-ayodeji/Linkrr/utils"
)

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
	html := utils.RenderTemplate("internal/email_templates/password_reset.html", data)
	sendEmail("Password-Reset", subject, text, html, name, email)
}
