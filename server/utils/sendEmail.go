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

// send email to user to change password or verify its account
func SendEmail(user *models.User, emailType string) bool {
	client, ctx := getMailSlurpClient()
	inbox, _, _ := client.InboxControllerApi.CreateInbox(ctx, nil)

	var body string
	var subject string
	var token string

	if emailType == "change" {
		token, _ = GenerateToken(user.ID, "email_confirmation")
		body = fmt.Sprintf(`<a href="%s/change-password/%s">reset password</a>`, os.Getenv("CORS_ORIGIN"), token)
		subject = "change password"

	} else if emailType == "verify" {
		token, _ = GenerateToken(user.ID, "change_password")
		body = fmt.Sprintf(`<a href="%s/change-password/%s">reset password</a>`, os.Getenv("CORS_ORIGIN"), token)
		subject = "verify account"
	}

	sendEmailOptions := MailSlurpClient.SendEmailOptions{
		To:      []string{inbox.EmailAddress},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}

	opts := &MailSlurpClient.SendEmailOpts{
		SendEmailOptions: optional.NewInterface(sendEmailOptions),
	}

	res, err := client.InboxControllerApi.SendEmail(ctx, inbox.Id, opts)
	if err != nil {
		return false
	}

	fmt.Println(res.StatusCode)
	return true
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
