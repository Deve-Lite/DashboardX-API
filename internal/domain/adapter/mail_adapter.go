package adapter

type MailAdapter interface {
	SendConfirmAccount(receiver string, token string) error
	SendPasswordReset(receiver string, token string) error
}
