package email

import (
	"fmt"

	"github.com/john-ayodeji/Linkrr/internal/utils"
)

func SendPasswordChangedEmail(name, email string) {
	subject := fmt.Sprintf("Password Change Successful")
	text := fmt.Sprintf("Hello %v\n your password has been changed successfully\nIf this action wasn't carried out by you, reset your password now.", name)
	data := struct {
		Name string
	}{
		Name: name,
	}
	html := utils.RenderTemplate("internal/email_templates/password_changed.html", data)
	sendEmail("Password-Reset", subject, text, html, name, email)
}
