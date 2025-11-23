package commands

import (
	"fmt"
	"net/smtp"
	"os"

	"koenbot/src/libs"
)

func checkSMTPActive(host, port, user, pass string) bool {
	from := user
	to := []string{os.Getenv("Email_Test")}

	// Email content
	subject := "Koenchan Test Email"
	body := "This is a test email sent"

	// Build the email message
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s",
		from,
		to[0], // For simplicity, assume a single recipient in the To header
		subject,
		body,
	))

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", user, pass, host)

	// Send the email
	err := smtp.SendMail(host+":"+port, auth, from, to, message)
	if err != nil {
		fmt.Printf("Error sending email: %v", err)
	}

	return err == nil
}

func splitLines(s string) []string {
	lines := []string{}
	curr := ""
	for _, r := range s {
		if r == '\n' || r == '\r' {
			if curr != "" {
				lines = append(lines, curr)
				curr = ""
			}
			continue
		}
		curr += string(r)
	}
	if curr != "" {
		lines = append(lines, curr)
	}
	return lines
}

func split(s, sep string) []string {
	var res []string
	curr := ""
	for i := 0; i < len(s); i++ {
		if len(s)-i >= len(sep) && s[i:i+len(sep)] == sep {
			res = append(res, curr)
			curr = ""
			i += len(sep) - 1
		} else {
			curr += string(s[i])
		}
	}
	res = append(res, curr)
	return res
}

func trim(s string) string {
	b := 0
	e := len(s)
	// Left trim
	for b < e && (s[b] == ' ' || s[b] == '\t' || s[b] == '\r' || s[b] == '\n') {
		b++
	}
	// Right trim
	for e > b && (s[e-1] == ' ' || s[e-1] == '\t' || s[e-1] == '\r' || s[e-1] == '\n') {
		e--
	}
	return s[b:e]
}

func join(arr []string, sep string) string {
	res := ""
	for i, s := range arr {
		if i > 0 {
			res += sep
		}
		res += s
	}
	return res
}

func writeFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func init() {
	libs.NewCommands(&libs.ICommand{
		Name: "(getsmtp|smtpc)",
		As:   []string{"getsmtp"},
		Description: "Mass check usable smtp from your list. Commands available: getsmtp,smtpc.\n\n" +
			"# Example:\n" +
			"```" +
			"getsmtp mail.google.com|45|kuntul@gmail.com|beraxsilied1337" +
			"mail.office.com|567|ffex1.3@office.com|long@or3#taEk" +
			"```" +
			"\n\nUp to 20 lines checking",
		Tags:     "hacking",
		IsPrefix: true,
		IsQuerry: true,
		IsWaitt:  true,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			// get list of smtp from input string (one per line)
			smtps := m.Querry
			lines := splitLines(smtps)
			// Limit to maximum 20 lines
			if len(lines) > 20 {
				lines = lines[:20]
			}

			if len(lines) == 0 {
				m.Reply("No SMTP data provided")
				return
			}

			activeSmtps := []string{}
			for _, smtpLine := range lines {
				smtpLine = trim(smtpLine)
				if smtpLine == "" {
					continue
				}
				parts := split(smtpLine, "|")
				if len(parts) < 4 {
					continue
				}
				host, port, user, pass := parts[0], parts[1], parts[2], parts[3]
				ok := checkSMTPActive(host, port, user, pass)
				if ok {
					activeSmtps = append(activeSmtps, smtpLine)
				}
			}
			if len(activeSmtps) == 0 {
				m.Reply("Sorry.. No active SMTPs found for you " + m.PushName + "-kun :(")
			} else {
				// Create sally-smtps.txt and send it
				filename := "sally-smtps.txt"
				content := join(activeSmtps, "\n")
				err := writeFile(filename, content)
				if err != nil {
					m.Reply("Failed to create file for you")
				} else {
					// Read file and send it
					fileData, err := os.ReadFile(filename)
					if err != nil {
						m.Reply("Failed to read file")
					} else {
						client.SendDocument(m.From, fileData, filename, "Here is your active SMTPs "+m.PushName+"-kun >//<", m.ID)
						// Clean up: delete the file after sending
						os.Remove(filename)
					}
				}
			}
		},
	})
}
