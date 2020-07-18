package sms

import (
	"time"
)

type Message struct {
	ID       string
	Time     time.Time
	From     string
	To       string
	Text     string
	Provider string
}
