package sms

type SmsNotifier interface {
	SmsNotify(msg Message)
}
