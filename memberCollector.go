package main

import (
	"log"

	"github.com/slack-go/slack"
)

// MemberCollector decide member randomly and make slackAttachment.
type MemberCollector struct {
	client *slack.Client
}

// Collect channnel members using slack api.
func (c *MemberCollector) Collect(channelID string) ([]Member, error) {
	var members []Member
	chInfo, err := c.client.GetChannelInfo(channelID)
	// err return when channel is private. Try GetGroupInfo.
	if err != nil {
		grInfo, err := c.client.GetGroupInfo(channelID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for _, mem := range grInfo.Members {
			user, _ := c.client.GetUserInfo(mem)
			if !(user.IsBot || user.Deleted) {
				members = append(members, Member{ID: user.ID, Name: user.Name})
			}
		}
	} else {
		for _, mem := range chInfo.Members {
			user, err := c.client.GetUserInfo(mem)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if !(user.IsBot || user.Deleted) {
				members = append(members, Member{ID: user.ID, Name: user.Name})
			}
		}
	}

	return members, nil
}
