package localmail

import (
	"encoding/base64"
	"net/smtp"
	"strings"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
)

func sendFromLocal(to []string, subject, body string) error {
	from := util.GetConfig("general", "send_from")
	addr := util.GetConfig("general", "mail_host")

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.Mail(from); err != nil {
		return err
	}
	for _, toOne := range to {
		if err = client.Rcpt(toOne); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	message := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

func init() {
	log.Info("Registering mail handler: local")
	util.HandleSendMail(func(mail util.MailMessage) error {
		return sendFromLocal(mail.To, mail.Subject, mail.Body)
	})
}
