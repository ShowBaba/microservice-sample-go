package app

type EmailMsgPayload struct {
	To      []string
	Subject string
	Body    string
}
