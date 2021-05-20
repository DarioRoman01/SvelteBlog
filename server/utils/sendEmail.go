package utils

import (
	"blogv2/users/models"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/antihax/optional"
	"github.com/joho/godotenv"
	MailSlurpClient "github.com/mailslurp/mailslurp-client-go"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("cannot read env")
	}
}

func SendChangePasswordEmail(user *models.User) bool {
	token, _ := GenerateToken(user.ID, CHANGE)
	body := fmt.Sprintf(`<a href="%s/change-password/%s">reset password</a>`, os.Getenv("CORS_ORIGIN"), token)
	subject := "change password"
	return sendEmail(token, body, subject)
}

func SendVerificationEmail(user *models.User) bool {
	token, _ := GenerateToken(user.ID, VERIFY)
	body := fmt.Sprintf(`<a href="%s/verify/%s">verify account</a>`, os.Getenv("CORS_ORIGIN"), token)
	subject := "verify account"
	return sendEmail(token, body, subject)
}

func getMailSlurpClient() (*MailSlurpClient.APIClient, context.Context) {
	ctx := context.WithValue(
		context.Background(),
		MailSlurpClient.ContextAPIKey,
		MailSlurpClient.APIKey{Key: os.Getenv("EMAIL_KEY")},
	)

	config := MailSlurpClient.NewConfiguration()
	client := MailSlurpClient.NewAPIClient(config)

	return client, ctx
}

func sendEmail(token, body, subject string) bool {
	client, ctx := getMailSlurpClient()
	inbox, _, _ := client.InboxControllerApi.CreateInbox(ctx, nil)

	sendEmailOptions := MailSlurpClient.SendEmailOptions{
		To:      []string{inbox.EmailAddress},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}

	opts := &MailSlurpClient.SendEmailOpts{
		SendEmailOptions: optional.NewInterface(sendEmailOptions),
	}

	_, err := client.InboxControllerApi.SendEmail(ctx, inbox.Id, opts)

	return err != nil
}
