package email

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"github.com/BattleBas/go-surprise/pkg/matching"
)

// Send out an email with the name of the reciever
func Send(m matching.Pair) error {

	from := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PASS")

	body := "Hello " + m.Giver.Name + ", \n You have been matched with " + m.Reciever.Name

	msg := "From: " + from + "\n" +
		"To: " + m.Giver.Email + "\n" +
		"Subject: Surprise\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{m.Giver.Email}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

// SendMasterList will send all the matches to the "surprise" email
func SendMasterList(m matching.Matches) error {

	from := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PASS")

	var body string

	for _, p := range m.Pairs {
		body += fmt.Sprintf("%s -> %s \n", p.Giver.Name, p.Reciever.Name)
	}

	msg := "From: " + from + "\n" +
		"To: " + from + "\n" +
		"Subject: Surprise Master List\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{from}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func IsValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
