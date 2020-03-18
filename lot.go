package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/slack-go/slack"
)

// Lot decide member randomly and make slackAttachment.
type Lot struct {
	client          *slack.Client
	messageTemplate MessageTemplate
}

// DrawLots decide member randomly and send message to slack.
func (l *Lot) DrawLots(channelID string, members []Member) error {
	if len(members) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	winner := members[rand.Intn(len(members))]
	attachment := slack.Attachment{
		Text:       fmt.Sprintf(l.messageTemplate.Choose, winner.Name),
		Color:      "#42f46e",
		CallbackID: "taskuji",
		Actions: []slack.AttachmentAction{
			{
				Name:  actionAccept,
				Text:  "OK!",
				Type:  "button",
				Style: "primary",
				Value: winner.ID,
			},
			{
				Name:  actionRepeat,
				Text:  "NG:cry:",
				Type:  "button",
				Style: "danger",
			},
		},
		Title:  l.messageTemplate.LotTitle,
		Footer: "Push the Button",
	}
	textMsg := slack.MsgOptionText(fmt.Sprintf("<@%s>", winner.ID), false)
	attachmentMsg := slack.MsgOptionAttachments(attachment)
	if _, _, err := l.client.PostMessage(channelID, textMsg, attachmentMsg); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}
