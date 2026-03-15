package email

import (
    "fmt"
    "net/smtp"
)

type Mailer struct {
    Host     string
    Port     string
    Username string
    Password string
    From     string
}

func NewMailer(host, port, username, password, from string) *Mailer {
    return &Mailer{
        Host:     host,
        Port:     port,
        Username: username,
        Password: password,
        From:     from,
    }
}

// Send sends an email
func (m *Mailer) Send(to, subject, body string) error {
    // Header formatting
    msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.From, to, subject, body)

    // Authentication
    auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

    // Sending
    addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
    err := smtp.SendMail(addr, auth, m.From, []string{to}, []byte(msg))
    
    // For debugging in console if Mailtrap fails
    if err != nil {
        fmt.Printf("Failed to send email to %s: %v\n", to, err)
        return err
    }
    
    fmt.Printf("Email sent successfully to %s (Check Mailtrap Inbox)\n", to)
    return nil
}