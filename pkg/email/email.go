package email

import (
	"net/smtp"
	"os"
)

type Match struct {
	Name     string
	Email    string
	Reciever string
}

// Send out an email with the name of the reciever
func Send(m Match) error {

	from := os.Getenv("EMAIL")
	pass := os.Getenv("PASSWORD")

	body := "Hello " + m.Name + ", \n You have been matched with " + m.Reciever

	msg := "From: " + from + "\n" +
		"To: " + m.Email + "\n" +
		"Subject: Surprise\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{m.Email}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
