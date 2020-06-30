package mailserver

type EmailNotifier interface {
	EmailNotify(msg MailMessage)
}
