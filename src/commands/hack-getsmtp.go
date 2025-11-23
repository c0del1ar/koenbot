package commands

import (
	"context"
	"net"
	"net/smtp"
	"time"

	"koenbot/src/libs"
)

func checkSMTPActive(host, port, user, pass string) bool {
	// Create context with timeout (10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create dialer with timeout
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	// Dial SMTP server
	addr := net.JoinHostPort(host, port)
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return false
	}
	defer client.Close()

	// Check if server supports AUTH
	authSupported, _ := client.Extension("AUTH")
	if !authSupported {
		return false
	}

	// Try to authenticate
	auth := smtp.PlainAuth("", user, pass, host)
	err = client.Auth(auth)
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

func init() {
	libs.NewCommands(&libs.ICommand{
		Name:        "(getsmtp|smtpc)",
		As:          []string{"getsmtp"},
		Description: "Check usable smtp from your list. Commands used: getsmtp,smtpc",
		Tags:        "hacking",
		IsPrefix:    true,
		IsQuerry:    true,
		IsWaitt:     true,
		Exec: func(client *libs.NewClientImpl, m *libs.IMessage) {
			// get list of smtp from input string (one per line)
			smtps := m.Querry
			lines := splitLines(smtps)

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
				m.Reply("No active SMTPs found.")
			} else {
				result := "Active SMTPs:\n" + join(activeSmtps, "\n")
				m.Reply(result)
			}
		},
	})
}
