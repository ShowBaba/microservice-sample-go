package shared

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO: replace all response payload with the appropriate struct below
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type LogPayload struct {
	Data   string
	Source string
}
type GatewayTokenJwtClaim struct {
	Gateway string
	jwt.StandardClaims
}

type AuthTokenJwtClaim struct {
	Email string
	jwt.StandardClaims
}

type APIResponse struct {
	Status  int      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func (mail *Mail) BuildMessage() string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	return msg
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Firstname string    `json:"firstname,omitempty"`
	Lastname  string    `json:"lastname,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}