package mail

import (
	"fmt"

	"testing"
)

func TestSendingMail(t *testing.T) {
	mailClient := MailClient{
		From:            "",
		ServiceAuthCode: "wkktdkdalndhbdia",
		ServiceDomain:   "smtp.qq.com",
		ServicePort:     "587",
	}

	if err := mailClient.SendEmail(
		[]string{}, "Hello! This is azusaings.",
		[]byte("Welcome to Revolution Fiesta"),
		[]byte("<h1>Welcome! This is Revolution Fiesta</h1><p>Congruatulations!</p><p>所以看起来中文不会乱码吧</p>"),
	); err != nil {
		fmt.Println(err)
	}
}
