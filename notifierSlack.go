package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"k8s.io/klog/v2"
)

type Slack struct {
	WebhookUrl     string
	DefaultChannel string // Slack channel name
	Username       string // Slack username (will show in notifier message)
}

var _ Notifier = &Slack{}

func newSlack() Notifier {
	var slackWebhookUrl, slackChannel, slackUsername string

	if slackWebhookUrl = os.Getenv("SLACK_WEBHOOK_URL"); slackWebhookUrl == "" {
		klog.Exit("Environment variable SLACK_WEBHOOK_URL is not set")
	}

	if slackChannel = os.Getenv("SLACK_CHANNEL"); slackChannel == "" {
		slackChannel = "restart-info-nonprod"
		klog.Warningf("Environment variable SLACK_CHANNEL is not set, default: %s\n", slackChannel)
	}

	if slackUsername = os.Getenv("SLACK_USERNAME"); slackUsername == "" {
		slackUsername = "k8s-pod-restart-info-collector"
		klog.Warningf("Environment variable SLACK_USERNAME is not set, default: %s\n", slackUsername)
	}

	klog.Infof("Slack Info: channel: %s, username: %s\n", slackChannel, slackUsername)

	return Slack{
		WebhookUrl:     slackWebhookUrl,
		DefaultChannel: slackChannel,
		Username:       slackUsername,
	}
}

func (s Slack) sendToChannel(msg Message, slackChannel string) error {
	channel := s.DefaultChannel
	if slackChannel != "" {
		channel = slackChannel
	}

	if len(msg.Text) > 8000 {
		// Slack attachment text will be truncated when > 8000 chars
		msg.Text = msg.Text[:7995] + "'''\n"
	}

	attachment := slack.Attachment{
		Text:       msg.Text,
		Title:      msg.Title,
		Footer:     msg.Footer,
		MarkdownIn: []string{"text"},
		Color:      "#4599DF",
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	err := slack.PostWebhook(s.WebhookUrl, &slack.WebhookMessage{
		Username:    s.Username,
		Channel:     channel,
		IconEmoji:   ":kubernetes:",
		Attachments: []slack.Attachment{attachment},
	})
	if err != nil {
		klog.Errorf("Sending to Slack channel failed with %v", err)
		return err
	}
	klog.Infof("Sent: [%s] to Slack.\n\n", strings.Replace(msg.Title, "\n", " ", -1))
	return nil
}
