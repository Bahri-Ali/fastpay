package email

import (
    "crypto/rand"
    "encoding/hex"
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

// Send sends an email with proper headers
func (m *Mailer) Send(to, subject, body string) error {
    // 1. Generate a random Message-ID (Required by RFC 5322)
    b := make([]byte, 16)
    rand.Read(b)
    messageID := fmt.Sprintf("<%s@%s>", hex.EncodeToString(b), m.Host)

    // 2. Construct Headers
    headers := make(map[string]string)
    headers["From"] = m.From
    headers["To"] = to
    headers["Subject"] = subject
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/plain; charset=\"utf-8\""
    headers["Message-ID"] = messageID

    // 3. Build Message String
    message := ""
    for k, v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body

    // 4. Authentication
    auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

    // 5. Send
    addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
    err := smtp.SendMail(addr, auth, m.From, []string{to}, []byte(message))

    if err != nil {
        fmt.Printf("Failed to send email to %s: %v\n", to, err)
        return err
    }

    fmt.Printf("Email sent successfully to %s (Check Mailtrap Inbox)\n", to)
    return nil
}