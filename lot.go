package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/seehuhn/mt19937"
	"github.com/slack-go/slack"
)

// Lot decide member randomly and make slackAttachment.
type Lot struct {
	client          *slack.Client
	messageTemplate MessageTemplate
}

// DrawLots decide member randomly and send message to slack.
func (l *Lot) DrawLots(channelID string, members []Member, userGroupID string) error {
	if len(members) == 0 {
		return errors.New("no member found")
	}
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	winner := draw(members, rng.Intn)
	messages := buildLotMessage(winner, userGroupID, l.messageTemplate)
	if _, _, err := l.client.PostMessage(channelID, messages...); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}

func draw(members []Member, rngfn func(int) int) Member {
	return members[rngfn(len(members))]
}

func buildLotMessage(winner Member, userGroupID string, template MessageTemplate) []slack.MsgOption {
	var messages []slack.MsgOption
	messages = append(messages, slack.MsgOptionText(fmt.Sprintf("<@%s>", winner.ID), false))

	attachment := slack.Attachment{
		Text:       fmt.Sprintf(template.Choose, winner.Name),
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
				Value: userGroupID,
			},
		},
		Title:  template.LotTitle,
		Footer: "Push the Button",
	}
	messages = append(messages, slack.MsgOptionAttachments(attachment))
	return messages
}
