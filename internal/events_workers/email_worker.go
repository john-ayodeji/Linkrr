package events_workers

import (
	"log"

	"github.com/john-ayodeji/Linkrr/internal/services/auth"
	"github.com/john-ayodeji/Linkrr/internal/services/email"
)

func SignUpEmailWorker(userData <-chan authService.UserData) {
	for data := range userData {
		email.SendWelcomeEmail(data.UserName, data.Email)
	}
}

func LoginEmailWorker(userData <-chan authService.UserData) {
	for data := range userData {
		email.SendLoginWelcomeEmail(data.UserName, data.Email)
	}
	log.Printf("login channel closed")
}

func ForgotPasswordEmailWorker(emailData <-chan authService.ForgotPasswordEmailData) {
	for data := range emailData {
		email.SendPasswordResetEmail(data.Name, data.Email, data.ResetURL)
	}
}

func ChangedPasswordEmailWorker(emailData <-chan authService.ResetPasswordEmailData) {
	for data := range emailData {
		email.SendPasswordChangedEmail(data.Name, data.Email)
	}
}
