package main

import "os"

const (
	NotifierTypeSlack = "notifier"
	NotifierTypeGChat = "gchat"
)

type Message struct {
	Title  string
	Text   string
	Footer string
}

type Notifier interface {
	sendToChannel(message Message, channel string) error
}

// GetNotifier based on environment variable "NOTIFIER_TYPE"
func GetNotifier() Notifier {
	switch os.Getenv("NOTIFIER") {
	case NotifierTypeSlack:
		return newSlack()
	default:
		return newSlack()
	}
}
